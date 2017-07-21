package main

import "debug/pe"
import "strconv"
import "os/exec"
import "runtime"
import "fmt"
import "os"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set the default values..
	peid.keySize = 8
	peid.staged = false
	peid.resource = true
	peid.verbose = false
	peid.iat = false

	ARGS := os.Args[1:]

	if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h" {
		Banner()
		Help()
		os.Exit(0)
	}

	Banner()
	peid.fileName = ARGS[0]

	for i := 0; i < len(ARGS); i++ {
		if ARGS[i] == "-ks" || ARGS[i] == "--keysize" {
			ks, Err := strconv.Atoi(ARGS[i+1])
			if Err != nil {
				boldRed.Println("\n[!] ERROR: Invalid key size.\n")
				fmt.Println(Err)
				os.Exit(1)
			} else {
				peid.keySize = ks
			}
		}
		if ARGS[i] == "-k" || ARGS[i] == "--key" {
			peid.key = []byte(ARGS[i+1])
		}
		if ARGS[i] == "--staged" {
			peid.staged = true
		}
		if ARGS[i] == "--no-resource" {
			peid.resource = false
		}
		if ARGS[i] == "-v" || ARGS[i] == "--verbose" {
			peid.verbose = true
		}
	}

	boldYellow.Print("\n[*] File: ")
	boldBlue.Println(peid.fileName)
	boldYellow.Print("[*] Staged: ")
	boldBlue.Println(peid.staged)
	if len(peid.key) != 0 {
		boldYellow.Print("[*] Key: ")
		boldBlue.Println(peid.key)
	} else {
		boldYellow.Print("[*] Key Size: ")
		boldBlue.Println(peid.keySize)
	}
	boldYellow.Print("[*] IAT: ")
	boldBlue.Println(peid.iat)
	boldYellow.Print("[*] Verbose: ")
	boldBlue.Println(peid.verbose, "\n")

	createBar()
	checkRequired() // 8 steps

	file, fileErr := pe.Open(ARGS[0])
	if fileErr != nil {
		boldRed.Println("\n[!] ERROR: Can't open file.")
		boldRed.Println(fileErr)
		os.Exit(1)
	}
	progress()

	analyze(file) // 9 steps
	assemble()    // 8 steps

	if peid.staged == true {
		exec.Command("sh", "-c", string("mv Payload "+peid.fileName+".stage")).Run()
	} else {
		compile() // 10 steps
	}
	clean() // 8 steps

	if peid.verbose == false {
		progressBar.Finish()
	}
	var getSize string = string("wc -c " + peid.fileName + "|tr -d \"" + peid.fileName + "\"|tr -d \"\n\"")

	if peid.staged == true {
		getSize = string("wc -c " + peid.fileName + "|tr -d \"" + peid.fileName + ".stage\"|tr -d \"\n\"")
	}

	wc, wcErr := exec.Command("sh", "-c", getSize).Output()
	if wcErr != nil {
		boldRed.Println("\n[!] ERROR: While getting the file size")
		boldRed.Println(string(wc))
		fmt.Println(wcErr)
		clean()
		os.Exit(1)
	}

	boldYellow.Println("\n[*] Final Size: " + peid.fileSize + "-> " + string(wc) + "bytes")
	if peid.staged == true {
		boldGreen.Println("[+] Stage generated !\n")
	} else {
		boldGreen.Println("[+] File successfully crypted !\n")
	}

}

func Banner() {

	var BANNER string = `

//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  POC Packer For Ophio                                             
`
	boldRed.Print(BANNER)
	boldBlue.Print("\n# Version: ")
	boldGreen.Println(VERSION)
	boldBlue.Print("# Source: ")
	boldGreen.Println("github.com/EgeBalci/Amber")

}

func Help() {
	var Help string = `

USAGE: 
  amber file.exe [options]


OPTIONS:
  
  -k, --key       [string]        Custom cipher key
  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:100/Min:4)
  --staged                        Generated a staged payload
  --no-resource                   Don't add any resource
  -v, --verbose                   Verbose output mode
  -h, --help                      Show this massage

EXAMPLE:
  (Default settings if no option parameter passed)
  amber file.exe -ks 8
`
	green.Println(Help)

}

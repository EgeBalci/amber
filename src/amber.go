package main

import "debug/pe"
import "strconv"
import "os/exec"
import "runtime"
import "os"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ARGS := os.Args[1:]
	if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h" {
		Banner()
		Help()
		os.Exit(0)
	}
	Banner()

	// Set the default values...
	peid.fileName = ARGS[0]
	peid.keySize = 7
	peid.staged = false
	peid.resource = true
	peid.verbose = false
	peid.iat = false

	// Parse the parameters...
	for i := 0; i < len(ARGS); i++ {
		if ARGS[i] == "-ks" || ARGS[i] == "--keysize" {
			ks, Err := strconv.Atoi(ARGS[i+1])
			if Err != nil {
				ParseError(Err,"\n[!] ERROR: Invalid key size.\n","")
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
		if ARGS[i] == "--iat" {
			peid.iat = true
		}
		if ARGS[i] == "--no-resource" {
			peid.resource = false
		}
		if ARGS[i] == "-v" || ARGS[i] == "--verbose" {
			peid.verbose = true
		}
	}
	// Show status
	BoldYellow.Print("\n[*] File: ")
	BoldBlue.Println(peid.fileName)
	BoldYellow.Print("[*] Staged: ")
	BoldBlue.Println(peid.staged)
	if len(peid.key) != 0 {
		BoldYellow.Print("[*] Key: ")
		BoldBlue.Println(peid.key)
	} else {
		BoldYellow.Print("[*] Key Size: ")
		BoldBlue.Println(peid.keySize)
	}
	BoldYellow.Print("[*] IAT: ")
	BoldBlue.Println(peid.iat)
	BoldYellow.Print("[*] Verbose: ")
	BoldBlue.Println(peid.verbose, "\n")

	// Create the process bar
	CreateProgressBar()
	CheckRequirements() // Check the required setup (6 steps)
	// Open the input file
	verbose("[*] Opening input file...",BoldYellow)
	file, fileErr := pe.Open(ARGS[0])
	ParseError(fileErr,"\n[!] ERROR: Can not open input file.","")
	progress()
	// Analyze the input file
	analyze(file) // 10 steps
	// Assemble the core amber payload
	assemble()    // 10 steps

	if peid.staged == true {
		exec.Command("sh", "-c", string("mv Payload "+peid.fileName+".stage")).Run()
	} else {
		compile() // Compile the amber stub (10 steps)
	}
	// Clean the created files
	clean() // 8 steps
	if peid.verbose == false {
		progressBar.Finish()
	}

	var getSize string = string("wc -c " + peid.fileName + "|tr -d \"" + peid.fileName + "\"|tr -d \"\n\"")
	if peid.staged == true {
		getSize = string("wc -c " + peid.fileName+ ".stage" + "|tr -d \"" + peid.fileName + ".stage\"|tr -d \"\n\"")
	}
	wc, wcErr := exec.Command("sh", "-c", getSize).Output()
	ParseError(wcErr,"\n[!] ERROR: While getting the file size",string(wc))

	BoldYellow.Println("\n[*] Final Size: " + peid.fileSize + "-> " + string(wc) + "bytes")
	if peid.staged == true {
		BoldGreen.Println("[+] Stage generated !\n")
	} else {
		BoldGreen.Println("[+] File successfully packed !\n")
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
//  POC Reflective PE Packer                                             
`
	BoldRed.Print(BANNER)
	BoldBlue.Print("\n# Version: ")
	BoldGreen.Println(VERSION)
	BoldBlue.Print("# Source: ")
	BoldGreen.Println("github.com/egebalci/Amber")

}

func Help() {
	var Help string = `

USAGE: 
  amber file.exe [options]


OPTIONS:
  
  -k, --key       [string]        Custom cipher key
  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:100/Min:4)
  --staged                        Generated a staged payload
  --iat                           Uses import address table entries instead of hash api
  --no-resource                   Don't add any resource
  -v, --verbose                   Verbose output mode
  -h, --help                      Show this massage

EXAMPLE:
  (Default settings if no option parameter passed)
  amber file.exe -ks 8
`
	green.Println(Help)

}

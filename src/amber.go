package main

import "path/filepath"
import "debug/pe"
import "strconv"
import "os/exec"
import "runtime"
import "strings"
import "os"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Run faster !

	ARGS := os.Args[1:]
	if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h" {
		Banner()
		Help()
		os.Exit(0)
	}
	Banner()

	// Set the default values...
	peid.KeySize = 8
	peid.staged = false
	peid.resource = true
	peid.verbose = false
	peid.iat = false

	// Parse the parameters...
	for i := 0; i < len(ARGS); i++ {
		if ARGS[i] == "-ks" || ARGS[i] == "--keysize" {
			ks, Err := strconv.Atoi(ARGS[i+1])
			if Err != nil {
				ParseError(Err,"Invalid key size.\n","")
			} else {
				peid.KeySize = ks
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
	// Get the absolute path of the file
	abs,abs_err := filepath.Abs(ARGS[0])
	ParseError(abs_err,"Can not open input file.","")
	peid.FileName = abs

	// Show status
	BoldYellow.Print("\n[*] File: ")
	BoldBlue.Println(peid.FileName)
	BoldYellow.Print("[*] Staged: ")
	BoldBlue.Println(peid.staged)
	if len(peid.key) != 0 {
		BoldYellow.Print("[*] Key: ")
		BoldBlue.Println(peid.key)
	} else {
		BoldYellow.Print("[*] Key Size: ")
		BoldBlue.Println(peid.KeySize)
	}
	BoldYellow.Print("[*] IAT: ")
	BoldBlue.Println(peid.iat)
	BoldYellow.Print("[*] Verbose: ")
	BoldBlue.Println(peid.verbose, "\n")

	// Create the process bar
	CreateProgressBar()
	CheckRequirements() // Check the required setup (6 steps)
	Cdir("/usr/share/Amber")
	progress()
	// Open the input file
	verbose("Opening input file...","*")
	file, FileErr := pe.Open(peid.FileName)
	ParseError(FileErr,"Can not open input file.","")
	progress()
	// Analyze the input file
	analyze(file) // 10 steps
	// Assemble the core amber payload
	assemble()    // 10 steps

	if peid.staged == true {
		//exec.Command("sh", "-c", string("mv Payload "+peid.FileName+".stage")).Run()
		move("Payload",string(peid.FileName+".stage"))
	} else {
		compile() // Compile the amber stub (10 steps)
	}
	// Clean the created files
	clean() // 10 steps
	if peid.verbose == false {
		progressBar.Finish()
	}

	var getSize string = string("wc -c " + peid.FileName + "|awk '{print $1}'|tr -d '\n'")
	if peid.staged == true {
		getSize = string("wc -c " + peid.FileName+ ".stage" + "|awk '{print $1}'|tr -d '\n'")
	}
	wc, wcErr := exec.Command("sh", "-c", getSize).Output()
	ParseError(wcErr,"While getting the file size",string(wc))

	BoldYellow.Print("\n[*] ")
	white.Println("Final Size: " + peid.fileSize + " -> " + string(wc) + " bytes")
	if peid.staged == true {
		BoldGreen.Println("[✔] Stage generated !\n")
	} else {
		BoldGreen.Println("[✔] File successfully packed !\n")
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
//  POC Reflective PE Packer ☣                                           
`

	var ArchBanner string = `


//    _____      _____ _______________________________ 
//   /  _  \    /     \\______   \_   _____/\______   \
//  /  /_\  \  /  \ /  \|    |  _/|    __)_  |       _/
// /    |    \/    Y    \    |   \|        \ |    |   \
// \____|__  /\____|__  /______  /_______  / |____|_  /
//         \/         \/       \/        \/         \/ 
// POC Reflective PE Packer
`


	DIST, _ := exec.Command("lsb_release", "-a").Output()
	if strings.Contains(string(DIST), "Arch") {
		PACKET_MANAGER = "pacman"
		BoldRed.Print(ArchBanner)
	}else{
		BoldRed.Print(BANNER)
	}


	//BoldRed.Print(BANNER)
	
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
  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:255/Min:5)
  --staged                        Generated a staged payload
  --iat                           Uses import address table entries instead of hash api
  --no-resource                   Don't add any resource
  -v, --verbose                   Verbose output mode
  -h, --help                      Show this massage

EXAMPLE:
  (Default settings if no option parameter passed)
  amber file.exe -ks 40
`
	green.Println(Help)

}

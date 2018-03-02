package main

import "path/filepath"
import "debug/pe"
import "strconv"
import "os/exec"
import "runtime"
import "errors"
import "os"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Run faster !

	ARGS := os.Args[1:]
	if len(ARGS) == 0 || ARGS[0] == "--help" || ARGS[0] == "-h" {
		Banner()
		BoldGreen.Println(Help)
		os.Exit(0)
	}
	Banner()

	// Set the default values...
	peid.KeySize = 0
	peid.staged = false
	peid.resource = true
	peid.verbose = false
	peid.iat = false
	peid.debug = false

	// Parse the parameters...
	for i := 0; i < len(ARGS); i++ {
		if ARGS[i] == "-ks" || ARGS[i] == "--keysize" {
			ks, Err := strconv.Atoi(ARGS[i+1])
			if Err != nil {
				ParseError(Err, "Invalid key size.\n", " ")
			} else {
				peid.KeySize = ks
			}
		}
		if ARGS[i] == "-k" || ARGS[i] == "--key" {
			peid.key = []byte(ARGS[i+1])
			if len([]byte(ARGS[i+1])) < 8 {
				ParseError(errors.New("Invalid key size !"), "Key size can't be smaller than 8 byte.\n", " ")
			}
			peid.KeySize = len([]byte(ARGS[i+1]))
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
		if ARGS[i] == "--debug" {
			peid.debug = true
		}
	}

	peid.FileName = ARGS[(len(ARGS) - 1)]

	if peid.KeySize == 0 && peid.staged == false {
		peid.KeySize = 8
	}

	// Show status
	BoldYellow.Print("\n[*] File: ")
	BoldBlue.Println(peid.FileName)
	BoldYellow.Print("[*] Staged: ")
	BoldBlue.Println(peid.staged)
	if len(peid.key) != 0 {
		BoldYellow.Print("[*] Key: ")
		BoldBlue.Println(string(peid.key))
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

	// Get the absolute path of the file
	abs, abs_err := filepath.Abs(ARGS[(len(ARGS) - 1)])
	ParseError(abs_err, "Can not open input file.", "")
	peid.FileName = abs
	progress()
	Cdir("/usr/share/Amber")
	progress()
	// Open the input file
	verbose("Opening input file...", "*")
	file, FileErr := pe.Open(peid.FileName)
	ParseError(FileErr, "Can not open input file.", "")
	progress()
	// Analyze the input file
	analyze(file) // 10 steps
	// Assemble the core amber payload
	assemble() // 10 steps

	if peid.staged == true {
		if peid.KeySize != 0 {
			crypt() // 4 steps
			Cdir("/usr/share/Amber/core")
			nasm, Err := exec.Command("nasm", "-f", "bin", "RC4.asm", "-o", "/usr/share/Amber/Payload").Output()
			ParseError(Err, "While assembling the RC4 decipher header.", string(nasm))
		}
		_copy("/usr/share/Amber/Payload", string(peid.FileName+".stage")) // Incase the file is on different volume
	} else {
		crypt()   // 4 steps
		compile() // Compile the amber stub (10 steps)
	}
	// Clean the created files
	clean() // 10 steps
	if peid.verbose == false {
		progressBar.Finish()
	}

	var getSize string = string("wc -c " + peid.FileName + "|awk '{print $1}'|tr -d '\n'")
	if peid.staged == true {
		getSize = string("wc -c " + peid.FileName + ".stage" + "|awk '{print $1}'|tr -d '\n'")
	}
	wc, wcErr := exec.Command("sh", "-c", getSize).Output()
	ParseError(wcErr, "While getting the file size", string(wc))

	BoldYellow.Print("\n[*] ")
	white.Println("Final Size: " + peid.fileSize + " -> " + string(wc) + " bytes")
	if peid.staged == true {
		BoldGreen.Println("[✔] Stage generated !\n")
	} else {
		BoldGreen.Println("[✔] File successfully packed !\n")
	}

}

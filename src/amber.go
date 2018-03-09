package main

import "path/filepath"
import "debug/pe"
import "os/exec"
import "runtime"
import "errors"
import "flag"
import "fmt"
import "os"	

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Run faster !

	flag.IntVar(&target.KeySize, "k", 0, "Size of the encryption key in bytes (Max:255/Min:8)")
	flag.IntVar(&target.KeySize, "keysize", 0, "Size of the encryption key in bytes (Max:255/Min:8)")
	flag.BoolVar(&target.staged, "r", false, "Generated a reflective payload")
	flag.BoolVar(&target.staged, "reflective", false, "Generated a reflective payload")
	flag.BoolVar(&target.resource, "no-resource", false, "Don't add any resource data")
	flag.BoolVar(&target.iat, "i", false, "Uses import address table entries instead of export address table")
	flag.BoolVar(&target.iat, "iat", false, "Uses import address table entries instead of export address table")
	flag.BoolVar(&target.verbose, "v", false, "Verbose output mode")
	flag.BoolVar(&target.verbose, "verbose", false, "Verbose output mode")
	flag.BoolVar(&target.debug, "debug", false, "Do not clean created files")
	flag.BoolVar(&target.IgnoreMappingSize, "ignore-mapping-size", false, "Ignore mapping size mismatch errors")
	flag.BoolVar(&target.IgnoreSectionAlignment, "ignore-section-alignment", false, "Ignore broken section alignment errors")
	flag.BoolVar(&target.help, "h", false, "Display this message")
	flag.BoolVar(&target.help, "H", false, "Display this message")
	flag.Parse()
	target.clean = false

	if len(os.Args) == 1 || target.help {
		Banner()
		fmt.Println(Help)
	 	os.Exit(0)
	}
	Banner()
	
	if target.KeySize == 0 && target.staged == false {
		target.KeySize = 8
	}

	if target.KeySize < 8 || target.KeySize > 255 {
		ParseError(errors.New("Key size must be between 8-255."), "Invalid key size.\n")
	}

	ARGS := flag.Args()
	if len(ARGS) == 0{
		ParseError(errors.New("Target file not specified."), "Target file not specified.\n")
	}
	target.FileName = ARGS[(len(ARGS)-1)]

	// Show status
	BoldYellow.Print("\n[*] File: ")
	BoldBlue.Println(target.FileName)
	BoldYellow.Print("[*] Staged: ")
	BoldBlue.Println(target.staged)
	BoldYellow.Print("[*] Key Size: ")
	BoldBlue.Println(target.KeySize)	
	BoldYellow.Print("[*] IAT: ")
	BoldBlue.Println(target.iat)
	BoldYellow.Print("[*] Verbose: ")
	BoldBlue.Println(target.verbose, "\n")

	// Create the process bar
	CreateProgressBar()
	CheckRequirements() // Check the required setup (6 steps)

	// Get the absolute path of the file
	abs, abs_err := filepath.Abs(ARGS[(len(ARGS) - 1)])
	ParseError(abs_err, "Can not open input file.")
	target.FileName = abs
	progress()
	Cdir("/usr/share/Amber")
	progress()
	// Open the input file
	verbose("Opening input file...", "*")
	file, FileErr := pe.Open(target.FileName)
	ParseError(FileErr, "Can not open input file.")
	progress()
	// Analyze the input file
	analyze(file) // 10 steps
	// Assemble the core amber payload
	assemble() // 10 steps

	if target.staged == true {
		if target.KeySize != 0 {
			crypt() // 4 steps
			Cdir("/usr/share/Amber/core")
			Err := exec.Command("nasm", "-f", "bin", "RC4.asm", "-o", "/usr/share/Amber/Payload").Run()
			ParseError(Err, "While assembling the RC4 decipher header.")
		}
		_copy("/usr/share/Amber/Payload", string(target.FileName+".stage")) // Incase the file is on different volume
	} else {
		crypt()   // 4 steps
		compile() // Compile the amber stub (10 steps)
	}
	// Clean the created files
	clean() // 10 steps
	if target.verbose == false {
		progressBar.Finish()
	}

	var getSize string = string("wc -c " + target.FileName + "|awk '{print $1}'|tr -d '\n'")
	if target.staged == true {
		getSize = string("wc -c " + target.FileName + ".stage" + "|awk '{print $1}'|tr -d '\n'")
	}
	wc, wcErr := exec.Command("sh", "-c", getSize).Output()
	ParseError(wcErr, "While getting the file size")

	BoldYellow.Print("\n[*] ")
	white.Println("Final Size: " + target.FileSize + " -> " + string(wc) + " bytes")
	if target.staged == true {
		BoldGreen.Println("[✔] Reflective stage generated !\n")
	} else {
		BoldGreen.Println("[✔] File successfully packed !\n")
	}

}
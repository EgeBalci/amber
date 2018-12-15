package main

import (
	"debug/pe"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Run faster !

	flag.IntVar(&target.KeySize, "k", 8, "Size of the encryption key in bytes (Max:255/Min:8)")
	flag.IntVar(&target.KeySize, "keysize", 8, "Size of the encryption key in bytes (Max:255/Min:8)")
	flag.BoolVar(&target.reflective, "r", false, "Generated a reflective payload")
	flag.BoolVar(&target.reflective, "reflective", false, "Generated a reflective payload")
	flag.BoolVar(&target.AntiAnalysis, "a", false, "Add anti-analysis functions inside stub (Binary size will increase)")
	flag.BoolVar(&target.AntiAnalysis, "anti-analysis", false, "Add anti-analysis functions inside stub (Binary size will increase)")
	flag.BoolVar(&target.resource, "no-resource", false, "Don't add any resource data")
	flag.BoolVar(&target.scrape, "s", false, "Scrape the PE header info (May break some files)")
	flag.BoolVar(&target.scrape, "scrape", false, "Scrape the PE header info (May break some files)")
	flag.BoolVar(&target.iat, "i", false, "Uses import address table entries instead of export address table")
	flag.BoolVar(&target.iat, "iat", false, "Uses import address table entries instead of export address table")
	flag.BoolVar(&target.verbose, "v", false, "Verbose output mode")
	flag.BoolVar(&target.verbose, "verbose", false, "Verbose output mode")
	flag.BoolVar(&target.debug, "debug", false, "Do not clean created files")
	flag.BoolVar(&target.IgnoreIntegrity, "ignore-integrity", false, "Ignore integrity check errors.")
	flag.BoolVar(&target.help, "h", false, "Display this message")
	flag.BoolVar(&target.help, "H", false, "Display this message")
	flag.Parse()

	banner()
	if len(os.Args) == 1 || target.help {
		fmt.Println(Help)
		os.Exit(0)
	}

	if target.KeySize < 8 || target.KeySize > 255 {
		parseErr(errors.New("invalid key size, key size must be between 8-255"))
	}

	ARGS := flag.Args()
	if len(ARGS) == 0 {
		parseErr(errors.New("target file not specified"))
	}
	target.FileName = ARGS[(len(ARGS) - 1)]

	// Show status
	printParams()
	// Create the progress bar
	createProgressBar()
	// Check dependencies
	requirements()
	// Get the absolute path of the file
	abs, err := filepath.Abs(ARGS[(len(ARGS) - 1)])
	parseErr(err)
	// Setup the working directory
	target.workdir = workdir(target.FileName)
	verbose("Setting up workdirectory at "+target.workdir, "*")
	mkdir(target.workdir)
	cdir(target.workdir)
	target.FileName = abs
	// Open the input file
	verbose("Opening input file...", "*")
	file, err := pe.Open(target.FileName)
	parseErr(err)
	analyze(file) // 10 steps
	// Assemble the core amber payload
	assemble() // 10 steps
	if target.reflective {
		copyFile(target.workdir+"/stage", target.FileName+".stage") // Incase the file is on different volume
	} else {
		compile() // Compile the amber stub (10 steps)
	}
	// Clean the created files
	clean() // 10 steps
	if !target.verbose {
		progressBar.Finish()
	}

	finalSize := ""

	if target.reflective == true {
		finalSize = fileSize(target.FileName + ".stage")
	} else {
		finalSize = fileSize(target.FileName)
	}

	BoldYellow.Print("\n[*] ")
	white.Println("Final Size: " + target.FileSize + " -> " + finalSize + " bytes")
	if target.reflective == true {
		BoldGreen.Println("[✔] Reflective stub generated !\n")
	} else {
		BoldGreen.Println("[✔] File successfully packed !\n")
	}

}

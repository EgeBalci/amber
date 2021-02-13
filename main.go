package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/EgeBalci/amber/config"
	amber "github.com/EgeBalci/amber/pkg"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Set globals...
var spinr = spinner.New(spinner.CharSets[9], 30*time.Millisecond)

func main() {

	banner()
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	// Configure the options from the flags/config file
	bp, encoder, err := config.ConfigureOptions(fs, os.Args[1:])
	if err != nil {
		config.PrintUsageErrorAndDie(err)
	}

	green := color.New(color.FgGreen).Add(color.Bold)
	spinr.Start()
	status("File: %s\n", bp.FileName)
	status("Build Stub: %t\n", bp.BuildStub)
	status("Encode Count: %d\n", encoder.EncodingCount)
	if bp.IAT {
		status("API: IAT\n")
	} else {
		status("API: CRC\n")
	}
	// First analyze PE and generate a blueprint
	spinr.Suffix = " Analyzing PE file..."
	eror(bp.Analyze())
	if !bp.HasRelocData {
		statusBad("%s has no relocation data.\n", bp.FileName)
		if bp.ImageBase != 0x400000 {
			statusBad("Can't switch to fixed address loader because ImageBase mismatch!\n")
		}
		status("Switching to fixed address loader...\n")
	}
	spinr.Suffix = " Assembling reflective payload..."
	payload, err := bp.AssemblePayload()
	eror(err)

	if encoder.EncodingCount > 0 {
		spinr.Suffix = " Encoding reflective payload..."
		encoder.SetArchitecture(bp.Architecture)
		payload, err = encoder.Encode(payload)
		eror(err)
	}

	if !bp.BuildStub {
		bp.FullFileName += ".bin"
	} else {
		// Construct EXE stub
		spinr.Suffix = " Building EXE stub..."
		payload, err = bp.CompileStub(payload)
		eror(err)
		bp.FullFileName = strings.ReplaceAll(bp.FullFileName, filepath.Ext(bp.FullFileName), "_packed.exe")
	}
	spinr.Stop()
	outFile, err := os.Create(bp.FullFileName)
	eror(err)
	outFile.Write(payload)
	defer outFile.Close()

	finSize, err := amber.GetFileSize(bp.FullFileName)
	eror(err)
	status("Final Size: %d bytes\n", finSize)
	status("Build File: %s\n", bp.FileName)
	green.Println("[✔] Reflective PE generated !")

}

// BANNER .
const BANNER string = `

//       █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//      ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//      ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//      ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//      ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//      ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  Reflective PE Packer ☣ Copyright (c) 2017 EGE BALCI
//      %s - %s

`

func banner() {
	green := color.New(color.FgGreen).Add(color.Bold)
	red := color.New(color.FgRed).Add(color.Bold)
	blue := color.New(color.FgBlue).Add(color.Bold)
	red.Printf(BANNER, green.Sprintf("v%s", amber.VERSION), blue.Sprintf("https://github.com/egebalci/amber"))
}

func status(formatstr string, a ...interface{}) {
	if spinr.Active() {
		spinr.Stop()
	}
	yellow := color.New(color.FgYellow).Add(color.Bold)
	yellow.Print("[*] ")
	fmt.Printf(formatstr, a...)
	spinr.Start()
}

func statusBad(formatstr string, a ...interface{}) {
	if spinr.Active() {
		spinr.Stop()
	}
	red := color.New(color.FgRed).Add(color.Bold)
	white := color.New(color.FgWhite).Add(color.Bold)
	red.Print("[!] ")
	white.Printf(formatstr, a...)
	spinr.Start()
}

func eror(err error) {
	if err != nil {
		if spinr.Active() {
			spinr.Stop()
		}
		red := color.New(color.FgRed).Add(color.Bold)
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			red.Print("[x] ")
			log.Fatalf("%s: %s\n", strings.ToUpper(strings.Split(details.Name(), ".")[1]), err)
		} else {
			red.Print("[x] ")
			log.Fatal(err)
		}
	}
}

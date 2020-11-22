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

	amber "github.com/EgeBalci/amber/pkg"
	sgn "github.com/EgeBalci/sgn/lib"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Set globals...
var spinr = spinner.New(spinner.CharSets[9], 30*time.Millisecond)

func main() {

	banner()
	bp := new(amber.Blueprint)
	encoder := sgn.NewEncoder()

	flag.StringVar(&bp.FileName, "f", "", "Input PE file")
	flag.BoolVar(&bp.IAT, "iat", false, "Use IAT API resolver block instead of CRC API resolver block")
	flag.BoolVar(&bp.IgnoreIntegrity, "ignore-checks", false, "Ignore integrity check errors.")
	flag.StringVar(&bp.CustomStubName, "stub", "", "Use custom stub file (experimental)")
	flag.IntVar(&encoder.ObfuscationLimit, "max", 5, "Maximum number of bytes for obfuscation")
	flag.IntVar(&encoder.EncodingCount, "e", 1, "Number of times to encode the generated reflective payload")
	buildStub := flag.Bool("build", false, "Build EXE stub that executes the generated reflective payload")

	green := color.New(color.FgGreen).Add(color.Bold)
	flag.Parse()

	if bp.FileName == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	spinr.Start()
	status("File: %s\n", bp.FileName)
	status("Build Stub: %t\n", *buildStub)
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

	if !*buildStub {
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

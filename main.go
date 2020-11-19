package main

import (
	amber "amber/lib"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	sgn "github.com/egebalci/sgn/lib"
	"github.com/fatih/color"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var bp amber.Blueprint
var red = color.New(color.FgRed).Add(color.Bold)
var blue = color.New(color.FgBlue).Add(color.Bold)
var green = color.New(color.FgGreen).Add(color.Bold)
var yellow = color.New(color.FgYellow).Add(color.Bold)
var spinr = spinner.New(spinner.CharSets[9], 50*time.Millisecond)

func main() {

	banner()
	file := flag.String("f", "", "Input PE file")
	encount := flag.Int("e", 1, "Number of times to encode the binary (increases overall size)")
	buildStub := flag.Bool("build", false, "Build EXE stub for executing the reflective payload")
	iat := flag.Bool("iat", false, "Use IAT API resolver block instead of CRC API resolver")
	ignoreIntegrity := flag.Bool("ignore-checks", false, "Ignore integrity check errors.")
	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Get the absolute path of the file
	abs, err := filepath.Abs(*file)
	eror(err)

	spinr.Start()
	status("File: %s\n", abs)
	status("Build Stub: %t\n", *buildStub)
	status("Encode Count: %d\n", *encount)
	if bp.IAT {
		status("API: IAT\n")
	} else {
		status("API: CRC\n")
	}

	spinr.Suffix = " Analyzing PE file..."
	bp, err := amber.Analyze(abs)
	eror(err)
	status("Relocation Data: %t\n", bp.HasRelocData)
	bp.EncodeCount = *encount
	bp.IAT = *iat
	bp.IgnoreIntegrity = *ignoreIntegrity

	spinr.Suffix = " Assembling reflective payload..."
	payload, err := bp.AssemblePayload()
	eror(err)

	if *encount > 0 {
		spinr.Suffix = " Encoding reflective payload..."
		// Create a new SGN encoder
		encoder := sgn.NewEncoder()
		encoder.SetArchitecture(bp.Architecture)
		encoder.EncodingCount = *encount
		payload, err = encoder.Encode(payload)
		eror(err)
	}

	if !*buildStub {
		bp.FileName += ".bin"
	} else {
		// Construct EXE stub
		spinr.Suffix = " Building EXE stub..."
		payload, err = bp.CompileStub(payload)
		eror(err)
		bp.FileName = strings.ReplaceAll(bp.FileName, ".", "_packed.")
	}
	spinr.Stop()
	outFile, err := os.Create(bp.FileName)
	eror(err)
	outFile.Write(payload)
	defer outFile.Close()

	finSize, err := amber.GetFileSize(bp.FileName)
	eror(err)
	status("Final Size: %d bytes\n", finSize)
	green.Println("[✔] Reflective stub generated !\n")

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
	red.Printf(BANNER, green.Sprintf("v%s", amber.VERSION), blue.Sprintf("https://github.com/egebalci/amber"))
}

func status(formatstr string, a ...interface{}) {
	if spinr.Active() {
		spinr.Stop()
		status(formatstr, a...)
		spinr.Start()
		return
	}
	yellow.Print("[*] ")
	fmt.Printf(formatstr, a...)
}

func eror(err error) {
	if err != nil {
		if spinr.Active() {
			spinr.Stop()
		}
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

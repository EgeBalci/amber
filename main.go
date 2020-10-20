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

func main() {

	banner()
	encount := flag.Int("e", 1, "Number of times to encode the binary (increases overall size)")
	reflective := flag.Bool("r", false, "Generated a reflective payload")
	iat := flag.Bool("iat", false, "Use IAT API resolver block instead of CRC API resolver")
	//verbose := flag.Bool("v", false, "Verbose output mode")
	ignoreIntegrity := flag.Bool("ignore-integrity", false, "Ignore integrity check errors.")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	ARGS := flag.Args()
	// Get the absolute path of the file
	abs, err := filepath.Abs(ARGS[(len(ARGS) - 1)])
	eror(err)
	bp, err := amber.Analyze(abs)
	eror(err)

	status("File: %s\n", abs)
	status("Reflective: %t\n", *reflective)
	status("Encode Count: %d\n", *encount)
	if bp.IAT {
		status("API: IAT\n")
	} else {
		status("API: CRC\n")
	}

	bp.EncodeCount = *encount
	bp.Reflective = *reflective
	bp.IAT = *iat
	bp.IgnoreIntegrity = *ignoreIntegrity

	payload, err := bp.AssemblePayload()
	eror(err)

	// Create a new SGN encoder
	encoder := sgn.NewEncoder()
	encoder.SetArchitecture(bp.Architecture)

	for i := 0; i < bp.EncodeCount; i++ {
		// Encode the binary
		encodedPayload, err := encoder.Encode(payload)
		eror(err)
		payload = encodedPayload
	}

	if bp.Reflective {
		bp.FileName += ".bin"
		file, err := os.Create(bp.FileName)
		eror(err)
		file.Write(payload)
		defer file.Close()
	} else {
		backupFile(bp.FileName)
		file, err := os.Create(bp.FileName)
		eror(err)
		file.Write(payload)
		defer file.Close()
	}

	finSize, err := amber.GetFileSize(bp.FileName)
	eror(err)

	status("Final Size: %d bytes\n", finSize)
	if bp.Reflective {
		green.Println("[✔] Reflective stub generated !\n")
	} else {
		green.Println("[✔] File successfully packed !\n")
	}

}

func backupFile(fileName string) {
	err := os.Rename(fileName, fileName+".bak")
	eror(err)
}

// BANNER .
const BANNER string = `

//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  Reflective PE Packer ☣ by EGE BALCI %s
//  >> %s <<

`

func banner() {
	red.Printf(BANNER, green.Sprintf("v%s", amber.VERSION), blue.Sprintf("https://github.com/egebalci/amber"))
}

func status(formatstr string, a ...interface{}) {
	yellow.Print("[*] ")
	fmt.Printf(formatstr, a...)
}

func eror(err error) {
	if err != nil {
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

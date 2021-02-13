package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	amber "github.com/EgeBalci/amber/pkg"
	sgn "github.com/EgeBalci/sgn/pkg"
	"github.com/fatih/color"
)

var usageStr = `
Usage: amber [options]
Options:
    -f, --file <file>        Input PE file
    -s, --stub <file>        Use custom stub file (experimental)
    -m, --max  <int>         Maximum number of bytes for obfuscation
    -e,         <int>        Number of times to encode the generated reflective payload
    -b, --build              Build EXE stub that executes the generated reflective payload
    --iat,                   Use IAT API resolver block instead of CRC API resolver block
    --ignore-checks,         Ignore integrity check errors.
    -h,                      Show this message
`

// PrintUsageErrorAndDie ...
func PrintUsageErrorAndDie(err error) {
	color.Red(err.Error())
	fmt.Println(usageStr)
	os.Exit(1)
}

// PrintHelpAndDie ...
func PrintHelpAndDie() {
	fmt.Println(usageStr)
	os.Exit(0)
}

// ConfigureOptions accepts a flag set and augments it with agentgo-server
// specific flags. On success, an options structure is returned configured
// based on the selected flags.
func ConfigureOptions(fs *flag.FlagSet, args []string) (*amber.Blueprint, *sgn.Encoder, error) {

	// Create empty options
	bp := &amber.Blueprint{}
	encoder := sgn.NewEncoder()

	// Define flags
	help := fs.Bool("h", false, "Show help message")
	fs.StringVar(&bp.FileName, "f", "", "Input PE file")
	fs.StringVar(&bp.FileName, "file", "", "Input PE file")
	fs.BoolVar(&bp.IAT, "iat", false, "Use IAT API resolver block instead of CRC API resolver block")
	fs.BoolVar(&bp.IgnoreIntegrity, "ignore-checks", false, "Ignore integrity check errors.")
	fs.StringVar(&bp.CustomStubName, "s", "", "Use custom stub file (experimental)")
	fs.StringVar(&bp.CustomStubName, "stub", "", "Use custom stub file (experimental)")
	fs.IntVar(&encoder.ObfuscationLimit, "max", 5, "Maximum number of bytes for obfuscation")
	fs.IntVar(&encoder.EncodingCount, "e", 1, "Number of times to encode the generated reflective payload")
	fs.BoolVar(&bp.BuildStub, "b", false, "Build EXE stub that executes the generated reflective payload")
	fs.BoolVar(&bp.BuildStub, "build", false, "Build EXE stub that executes the generated reflective payload")

	// Parse arguments and check for errors
	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	// If it is not help and other args are empty, return error
	if (*help == false) && bp.FileName == "" {
		err := errors.New("please specify all required arguments")
		return nil, nil, err
	}

	// If -help flag is defined, print help
	if *help {
		PrintHelpAndDie()
	}

	return bp, encoder, nil
}

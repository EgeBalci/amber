package main

import (
	_ "github.com/egebalci/amber/statik"
	"github.com/fatih/color"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// VERSION number
const VERSION string = "2.0.0"

// ID structure for storing
// PE specs, tool parameters and
// OS spesific info globaly
type ID struct {

	// Parameters...
	FileName        string
	KeySize         int
	key             []byte
	reflective      bool
	AntiAnalysis    bool
	iat             bool
	resource        bool
	scrape          bool
	verbose         bool
	debug           bool
	help            bool
	clean           bool
	IgnoreIntegrity bool

	// System anaylsis
	nasm    string
	workdir string

	// PE analysis
	FileSize  string
	arch      string
	ImageBase uint64
	subsystem uint16
	aslr      bool
	dll       bool
}

var red *color.Color = color.New(color.FgRed)
var BoldRed *color.Color = red.Add(color.Bold)
var blue *color.Color = color.New(color.FgBlue)
var BoldBlue *color.Color = blue.Add(color.Bold)
var yellow *color.Color = color.New(color.FgYellow)
var BoldYellow *color.Color = yellow.Add(color.Bold)
var green *color.Color = color.New(color.FgGreen)
var BoldGreen *color.Color = green.Add(color.Bold)
var white *color.Color = color.New(color.FgWhite)

var progressBar *pb.ProgressBar
var target ID

// BANNER .
const BANNER string = `

//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  Reflective PE Packer ☣   
`

// BasicBanner .
const BasicBanner string = `


//    _____      _____ _______________________________ 
//   /  _  \    /     \\______   \_   _____/\______   \
//  /  /_\  \  /  \ /  \|    |  _/|    __)_  |       _/
// /    |    \/    Y    \    |   \|        \ |    |   \
// \____|__  /\____|__  /______  /_______  / |____|_  /
//         \/         \/       \/        \/         \/ 
// Reflective PE Packer
`

// Help .
const Help string = `
USAGE: 
  amber [options] file.exe
OPTIONS:
  -k, -keysize                Size of the encryption key in bytes (Max:255/Min:8)
  -r, -reflective             Generated a reflective payload
  -a, -anti-analysis          Add anti-analysis measures
  -i, -iat                    Use import address table entries instead of export address table
  -s, -scrape                 Scrape the PE header info (May break some files)
  -no-resource                Don't add any resource data (removes icon)
  -ignore-integrity           Ignore integrity check errors
  -v, -verbose                Verbose output mode
  -h, -H                      Show this massage
EXAMPLE:
  (Default settings if no option parameter passed)
  amber -k 8 file.exe
`

package main

import "gopkg.in/cheggaaa/pb.v1"
import "github.com/fatih/color"
import "debug/pe"

const VERSION string = "1.1.0"

var PACKET_MANAGER string = "apt"

type peID struct {

	// Parameters...
	FileName string
	KeySize  int
	key      []byte
	staged   bool
	iat      bool
	resource bool
	verbose  bool
	debug	 bool

	//Analysis...
	fileSize  string
	imageBase uint32
	subsystem uint16
	aslr      bool
	Opt       *pe.OptionalHeader32
	VP        string
	GPA       string
	LLA       string
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
//var BoldWhite *color.Color = white.Add(color.Bold)

var progressBar *pb.ProgressBar
var peid peID


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

	var Help string = `
USAGE: 
  amber [options] file.exe
OPTIONS:
  -k, --key            Custom cipher key
  -ks,--keysize        Size of the encryption key in bytes (Max:255/Min:8)
  --staged             Generated a staged payload
  --iat                Uses import address table entries instead of hash api
  --no-resource        Don't add any resource data
  -v, --verbose        Verbose output mode
  -h, --help           Show this massage
EXAMPLE:
  (Default settings if no option parameter passed)
  amber -ks 8 file.exe
`
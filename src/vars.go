package main

import "gopkg.in/cheggaaa/pb.v1"
import "github.com/fatih/color"
import "debug/pe"

const VERSION string = "1.3.0"

var PACKET_MANAGER string = "apt"

type PEID struct {

	// Parameters...
	FileName string
	KeySize  int
	key      []byte
	staged   bool
	iat      bool
	resource bool
	scrape	 bool
	verbose  bool
	debug    bool
	help 	 bool
	clean	 bool

	IgnoreIntegrity bool

	//Analysis...
	FileSize  string
	ImageBase uint32
	subsystem uint16
	aslr      bool
	Opt       *pe.OptionalHeader32
	VP        string
	GPA       string
	LLA       string
}

type OptionalHeader struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
	DataDirectory               [16]DataDirectory
}

type DataDirectory struct {
	VirtualAddress uint32
	Size           uint32
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
var target PEID

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
  -k, -keysize                Size of the encryption key in bytes (Max:255/Min:8)
  -r, -reflective             Generated a reflective payload
  -i, -iat                    Uses import address table entries instead of export address table
  -s, -scrape                 Scrape the PE header info (May break some files)
  -no-resource                Don't add any resource data
  -ignore-integrity           Ignore integrity check errors.
  -v, -verbose                Verbose output mode
  -h, -H                      Show this massage
EXAMPLE:
  (Default settings if no option parameter passed)
  amber -k 8 file.exe
`

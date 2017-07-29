package main

import "gopkg.in/cheggaaa/pb.v1"
import "github.com/fatih/color"
import "debug/pe"

const VERSION string = "1.0.0"

type peID struct {

	// Parameters...
	fileName string
	keySize  int
	key      []byte
	staged   bool
	iat      bool
	resource bool
	verbose  bool

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

var progressBar *pb.ProgressBar
var peid peID

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
	OPT       *pe.OptionalHeader32
	VP        string
	GPA       string
	LLA       string
}

var red *color.Color = color.New(color.FgRed)
var boldRed *color.Color = red.Add(color.Bold)
var blue *color.Color = color.New(color.FgBlue)
var boldBlue *color.Color = blue.Add(color.Bold)
var yellow *color.Color = color.New(color.FgYellow)
var boldYellow *color.Color = yellow.Add(color.Bold)
var green *color.Color = color.New(color.FgGreen)
var boldGreen *color.Color = green.Add(color.Bold)
var white *color.Color = color.New(color.FgWhite)

var progressBar *pb.ProgressBar
var peid peID

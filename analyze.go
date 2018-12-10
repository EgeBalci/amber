package main

import (
	"debug/pe"
	"errors"
	"fmt"

	"github.com/egebalci/mappe/mape"
)

func analyze(file *pe.File) {
	//Do analysis on pe file...
	verbose("Analyzing PE file...", "*")
	verbose("File Size: "+target.FileSize+" byte", "*")
	verbose("Machine: "+fmt.Sprintf("0x%X", file.FileHeader.Machine), "*")
	if file.FileHeader.Machine == 0x14C {
		target.arch = "x86"
	} else if file.FileHeader.Machine == 0x8664 {
		target.arch = "x64"
	} else {
		parseErr(errors.New("file architechture not supported"))
	}
	opt := mape.ConvertOptionalHeader(file)

	verbose("Magic: "+fmt.Sprintf("0x%X", opt.Magic), "*")
	if file.Characteristics >= 0x2000 {
		verbose("Found DLL characteristics.", "*")
		target.dll = true
	} else {
		verbose("Found executable characteristics.", "*")
		target.dll = false
	}
	if (opt.DataDirectory[5].Size) != 0x00 {
		target.aslr = true
		verbose("ASLR supported !", "+")
		verbose("Using ASLR stub...", "*")
	} else {
		target.aslr = false
		verbose("ASLR not supported :(", "-")
		verbose("Using Fixed stub...", "*")
	}
	if (opt.DataDirectory[11].Size) != 0x00 {
		parseErr(errors.New("file has bounded imports, EXE files with bounded imports not supported"))
	}
	if (opt.DataDirectory[14].Size) != 0x00 {
		parseErr(errors.New("Unempty CLR section, .NET binaries are not supported"))
	}
	if (opt.DataDirectory[13].Size) != 0x00 {
		verbose("WARNING: File has delayed imports. (This could be a problem :/ )", "!")
	}
	if (opt.DataDirectory[1].Size) == 0x00 {
		parseErr(errors.New("Import table size zero, file has empty import table"))
	}
	target.FileSize = fileSize(target.FileName)
	target.ImageBase = opt.ImageBase
	target.subsystem = opt.Subsystem

	verbose("Subsystem: "+fmt.Sprintf("0x%X", uint64(opt.Subsystem)), "*")
	verbose("Image Base: "+fmt.Sprintf("0x%X", uint64(opt.ImageBase)), "*")
	verbose("Address Of Entry: "+fmt.Sprintf("0x%X", uint64(opt.AddressOfEntryPoint)), "*")
	verbose("Size Of Image: "+fmt.Sprintf("0x%X", uint64(opt.SizeOfImage)), "*")
	verbose("Export Table: "+fmt.Sprintf("0x%X", uint64(opt.DataDirectory[0].VirtualAddress+uint32(opt.ImageBase))), "*")
	verbose("Import Table: "+fmt.Sprintf("0x%X", uint64(opt.DataDirectory[1].VirtualAddress+uint32(opt.ImageBase))), "*")
	verbose("Base Relocation Table: "+fmt.Sprintf("0x%X", uint64(opt.DataDirectory[5].VirtualAddress+uint32(opt.ImageBase))), "*")
	verbose("Import Address Table: "+fmt.Sprintf("0x%X", uint64(opt.DataDirectory[12].VirtualAddress+uint32(opt.ImageBase))), "*")

	if !target.aslr && target.ImageBase != 0x400000 {
		parseErr(errors.New("Unable to align image (file may already been packed)"))
	}

	defer progress()
}

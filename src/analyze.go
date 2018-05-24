package main

import "debug/pe"
import "os/exec"
import "errors"

func analyze(file *pe.File) {
	//Do analysis on pe file...
	verbose("Analyzing PE file...", "*")
	verbose("File Size: "+target.FileSize+" byte", "*")
	if file.FileHeader.Machine != 0x14C {
		ParseError(errors.New("False machine value !"), "File is not a 32 bit PE.")
	}
	progress()
	_verbose("Machine:", uint64(file.FileHeader.Machine))
	if file.Machine == 0x8664 {
		_opt := (file.OptionalHeader.(*pe.OptionalHeader64))
		target.opt.Magic = _opt.Magic
		target.opt.Subsystem = _opt.Subsystem
		target.opt.CheckSum = _opt.CheckSum
		target.opt.ImageBase = _opt.ImageBase
		target.opt.AddressOfEntryPoint = _opt.AddressOfEntryPoint
		target.opt.SizeOfImage =  _opt.SizeOfImage
		target.opt.SizeOfHeaders = _opt.SizeOfHeaders
		for i:=0; i<16; i++ {
			target.opt.DataDirectory[i].VirtualAddress = _opt.DataDirectory[i].VirtualAddress
			target.opt.DataDirectory[i].Size = _opt.DataDirectory[i].Size
		}
	}else{
		_opt := file.OptionalHeader.((*pe.OptionalHeader32))
		target.opt.Magic = _opt.Magic
		target.opt.Subsystem = _opt.Subsystem
		target.opt.CheckSum = _opt.CheckSum
		target.opt.ImageBase = uint64(_opt.ImageBase)
		target.opt.AddressOfEntryPoint = _opt.AddressOfEntryPoint
		target.opt.SizeOfImage =  _opt.SizeOfImage
		target.opt.SizeOfHeaders = _opt.SizeOfHeaders
		for i:=0; i<16; i++ {
			target.opt.DataDirectory[i].VirtualAddress = _opt.DataDirectory[i].VirtualAddress
			target.opt.DataDirectory[i].Size = _opt.DataDirectory[i].Size
		}
	}
	progress()
	_verbose("Magic:", uint64(target.opt.Magic))
	// PE32 = 0x10B
	if target.opt.Magic != 0x10B {
		ParseError(errors.New("False magic value !"), "File is not a valid PE.")
	}
	progress()
	if file.Characteristics >= 0x2000 {
		verbose("Found DLL characteristics.", "*")
		target.dll = true
	}else {
		verbose("Found executable characteristics.", "*")
		target.dll = false
	}
	progress()
	if (target.opt.DataDirectory[5].Size) != 0x00 {
		target.aslr = true
		verbose("ASLR supported !","+")
		verbose("Using ASLR stub...","*")	
	} else if (target.opt.DataDirectory[5].Size) == 0x00 {
		target.aslr = false
		verbose("ASLR not supported :(","-")
		verbose("Using Fixed stub...","*")
	}
	progress()
	if (target.opt.DataDirectory[11].Size) != 0x00 {
		ParseError(errors.New("File has bounded imports."), "EXE files with bounded imports not supported.")
	}
	progress()
	if (target.opt.DataDirectory[14].Size) != 0x00 {
		ParseError(errors.New("Unempty CLR section !"), ".NET binaries are not supported.")
	}
	progress()
	if (target.opt.DataDirectory[13].Size) != 0x00 {
		verbose("WARNING: File has delayed imports. (This could be a problem :/ )", "!")
	}
	progress()
	if (target.opt.DataDirectory[1].Size) == 0x00 {
		ParseError(errors.New("Import table size zero !"), "File has empty import table.")
	}
	progress()
	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+target.FileName+"|awk '{print $1}'|tr -d '\n'")).Output()
	ParseError(wcErr, "While getting the file size")
	target.FileSize = string(wc)
	progress()

	
	
	_verbose("Subsystem:", uint64(target.opt.Subsystem))
	_verbose("Image Base:", uint64(target.opt.ImageBase))
	_verbose("Address Of Entry:", uint64(target.opt.AddressOfEntryPoint))
	_verbose("Size Of Image:", uint64(target.opt.SizeOfImage))
	_verbose("Export Table:", uint64(target.opt.DataDirectory[0].VirtualAddress+uint32(target.opt.ImageBase)))
	_verbose("Import Table:", uint64(target.opt.DataDirectory[1].VirtualAddress+uint32(target.opt.ImageBase)))
	_verbose("Base Relocation Table:", uint64(target.opt.DataDirectory[5].VirtualAddress+uint32(target.opt.ImageBase)))
	_verbose("Import Address Table:", uint64(target.opt.DataDirectory[12].VirtualAddress+uint32(target.opt.ImageBase)))

}

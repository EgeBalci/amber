package main

import "debug/pe"
import "os/exec"
import "errors"

func analyze(file *pe.File) {
	//Do analysis on pe file...
	verbose("Analyzing PE file...", "*")
	if file.FileHeader.Machine != 0x14C {
		ParseError(errors.New("False machine value !"), "File is not a 32 bit PE.")
	}
	progress()
	var Opt *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if Opt.Magic != 0x10B {
		ParseError(errors.New("False magic value !"), "File is not a valid PE.")
	}
	progress()
	target.ImageBase = Opt.ImageBase
	progress()
	target.subsystem = Opt.Subsystem
	progress()
	if (Opt.DataDirectory[5].Size) != 0x00 {
		target.aslr = true
		if target.verbose == true {
			BoldGreen.Println("[+] ASLR supported !")
			BoldYellow.Println("[x] Using ASLR stub...")
		}
	} else if (Opt.DataDirectory[5].Size) == 0x00 {
		target.aslr = false
		if target.verbose == true {
			BoldRed.Println("[-] ASLR not supported :(")
			BoldYellow.Println("[x] Using Fixed stub...")
		}
	}
	progress()
	if (Opt.DataDirectory[11].Size) != 0x00 {
		ParseError(errors.New("File has bounded imports."), "EXE files with bounded imports not supported.")
	}
	progress()
	if (Opt.DataDirectory[14].Size) != 0x00 {
		ParseError(errors.New("Unempty CLR section !"), ".NET binaries are not supported.")
	}
	progress()
	if (Opt.DataDirectory[13].Size) != 0x00 {
		verbose("WARNING: File has delayed imports. (This could be a problem :/ )", "!")
	}
	progress()
	if (Opt.DataDirectory[1].Size) == 0x00 {
		ParseError(errors.New(""), "File has empty import table.")
	}
	progress()
	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+target.FileName+"|awk '{print $1}'|tr -d '\n'")).Output()
	ParseError(wcErr, "While getting the file size")
	target.FileSize = string(wc)
	progress()

	target.Opt = Opt

	verbose(string("File Size: "+target.FileSize+" byte"), "*")
	_verbose("Machine:", uint64(file.FileHeader.Machine))
	_verbose("Magic:", uint64(Opt.Magic))
	_verbose("Subsystem:", uint64(Opt.Subsystem))
	_verbose("Image Base:", uint64(target.ImageBase))
	_verbose("Address Of Entry:", uint64(Opt.AddressOfEntryPoint))
	_verbose("Size Of Image:", uint64(Opt.SizeOfImage))
	_verbose("Export Table:", uint64(Opt.DataDirectory[0].VirtualAddress+Opt.ImageBase))
	_verbose("Import Table:", uint64(Opt.DataDirectory[1].VirtualAddress+Opt.ImageBase))
	_verbose("Base Relocation Table:", uint64(Opt.DataDirectory[5].VirtualAddress+Opt.ImageBase))
	_verbose("Import Address Table:", uint64(Opt.DataDirectory[12].VirtualAddress+Opt.ImageBase))

}

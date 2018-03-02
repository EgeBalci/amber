package main

import "debug/pe"
import "os/exec"
import "errors"

func analyze(file *pe.File) {
	//Do analysis on pe file...
	verbose("Analyzing PE file...", "*")
	if file.FileHeader.Machine != 0x14C {
		ParseError(errors.New("False machine value !"), "File is not a 32 bit PE.", "")
	}
	progress()
	var Opt *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if Opt.Magic != 0x10B {
		ParseError(errors.New("False magic value !"), "File is not a valid PE.", "")
	}
	progress()
	peid.imageBase = Opt.ImageBase
	progress()
	peid.subsystem = Opt.Subsystem
	progress()
	if (Opt.DataDirectory[5].Size) != 0x00 {
		peid.aslr = true
		if peid.verbose == true {
			BoldGreen.Println("[+] ASLR supported !")
			BoldYellow.Println("[x] Using ASLR stub...")
		}
	} else if (Opt.DataDirectory[5].Size) == 0x00 {
		peid.aslr = false
		if peid.verbose == true {
			BoldRed.Println("[-] ASLR not supported :(")
			BoldYellow.Println("[x] Using Fixed stub...")
		}
	}
	progress()
	if (Opt.DataDirectory[11].Size) != 0x00 {
		ParseError(errors.New("File has bounded imports."), "EXE files with bounded imports not supported.", "")
	}
	progress()
	if (Opt.DataDirectory[14].Size) != 0x00 {
		ParseError(errors.New("Unempty CLR section !"), ".NET binaries are not supported.", "")
	}
	progress()
	if (Opt.DataDirectory[13].Size) != 0x00 {
		verbose("WARNING: File has delayed imports. (This could be a problem :/ )", "!")
	}
	progress()
	if (Opt.DataDirectory[1].Size) == 0x00 {
		ParseError(errors.New(""), "File has empty import table.", "")
	}
	progress()
	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+peid.FileName+"|awk '{print $1}'|tr -d '\n'")).Output()
	ParseError(wcErr, "While getting the file size", string(wc))
	peid.fileSize = string(wc)
	progress()

	peid.Opt = Opt

	verbose(string("File Size: "+peid.fileSize+" byte"), "*")
	_verbose("Machine:", int32(file.FileHeader.Machine))
	_verbose("Magic:", int32(Opt.Magic))
	_verbose("Subsystem:", int32(Opt.Subsystem))
	_verbose("Image Base:", int32(peid.imageBase))
	_verbose("Size Of Image:", int32(Opt.SizeOfImage))
	_verbose("Export Table:", int32(Opt.DataDirectory[0].VirtualAddress+Opt.ImageBase))
	_verbose("Import Table:", int32(Opt.DataDirectory[1].VirtualAddress+Opt.ImageBase))
	_verbose("Base Relocation Table:", int32(Opt.DataDirectory[5].VirtualAddress+Opt.ImageBase))
	_verbose("Import Address Table:", int32(Opt.DataDirectory[12].VirtualAddress+Opt.ImageBase))

}

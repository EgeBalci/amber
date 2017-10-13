package main

import "debug/pe"
import "os/exec"
import "errors"


func analyze(file *pe.File) {
	//Do analysis on pe file...
	verbose("[*] Analyzing PE file...",BoldYellow)
	if file.FileHeader.Machine != 0x14C {
		ParseError(errors.New(""),"\n[!] ERROR: File is not a 32 bit PE.","")
	}
	progress()
	var Opt *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if Opt.Magic != 0x10B {
		ParseError(errors.New(""),"\n[!] ERROR: File is not a valid PE.","")
	}
	progress()
	peid.imageBase = Opt.ImageBase
	progress()
	peid.subsystem = Opt.Subsystem
	progress()
	if (Opt.DataDirectory[5].Size) != 0x00 && peid.verbose == true {
		peid.aslr = true
		BoldGreen.Println("[+] ASLR supported !")
		BoldYellow.Println("[*] Using ASLR stub...")
	} else if (Opt.DataDirectory[5].Size) == 0x00 && peid.verbose == true {
		peid.aslr = false
		BoldRed.Println("[-] ASLR not supported :(")
		BoldYellow.Println("[*] Using Fixed stub...")
	}
	progress()
	if (Opt.DataDirectory[11].Size) != 0x00 {
		ParseError(errors.New(""),"\n[!] ERROR: File has bounded imports.","")
	}
	progress()
	if (Opt.DataDirectory[14].Size) != 0x00 {
		ParseError(errors.New(""),"\n[!] ERROR: .NET binaries are not supported.","")
	}
	progress()
	if (Opt.DataDirectory[13].Size) != 0x00 {
		verbose("[!] WARNING: File has delayed imports. (This could be a problem :/ )",BoldYellow)
	}
	progress()
	if (Opt.DataDirectory[1].Size) == 0x00 {
		ParseError(errors.New(""),"\n[!] ERROR: File has empty import table.","")
	}
	progress()
	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+peid.fileName+"|awk '{print $1}'|tr -d '\n'")).Output()
	ParseError(wcErr,"\n[!] ERROR: While getting the file size",string(wc))
	peid.fileSize = string(wc)
	progress()

	peid.Opt = Opt

	if peid.verbose == true {
		BoldYellow.Println("[*] File Size: " + peid.fileSize + " byte")
		BoldYellow.Printf("[*] Machine: 0x%X\n", file.FileHeader.Machine)
		BoldYellow.Printf("[*] Magic: 0x%X\n", Opt.Magic)
		BoldYellow.Printf("[*] Subsystem: 0x%X\n", Opt.Subsystem)
		BoldYellow.Printf("[*] Image Base: 0x%X\n", peid.imageBase)
		BoldYellow.Printf("[*] Size Of Image: 0x%X\n", Opt.SizeOfImage)
		BoldYellow.Printf("[*] Export Table: 0x%X\n", (Opt.DataDirectory[0].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Import Table: 0x%X\n", (Opt.DataDirectory[1].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Base Relocation Table: 0x%X\n", (Opt.DataDirectory[5].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Import Address Table: 0x%X\n", (Opt.DataDirectory[12].VirtualAddress + Opt.ImageBase))
	}
}

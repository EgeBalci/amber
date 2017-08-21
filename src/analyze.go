package main

import "debug/pe"
import "os/exec"
import "os"

func analyze(file *pe.File) {
	//Do analysis on pe file...

	if file.FileHeader.Machine != 0x14C {
		BoldRed.Println("\n[!] ERROR: File is not a 32 bit PE.")
		os.Exit(1)
	}
	progress()
	var Opt *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if Opt.Magic != 0x10B {
		BoldRed.Println("\n[!] ERROR: File is not a valid PE.")
		os.Exit(1)
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
			BoldYellow.Println("[*] Using ASLR stub...")
		}
	} else {
		peid.aslr = false
		if peid.verbose == true {
			BoldYellow.Println("[-] ASLR not supported :(")
			BoldYellow.Println("[*] Using Non-ASLR stub...")
		}
	}
	progress()

	if (Opt.DataDirectory[11].Size) != 0x00 {
		BoldRed.Println("\n[!] ERROR: File has bounded imports.")
		os.Exit(1)
	}
	progress()
	if (Opt.DataDirectory[13].Size) != 0x00 {
		BoldRed.Println("\n[!] ERROR: File has delayed imports.")
		os.Exit(1)
	}
	progress()

	if (Opt.DataDirectory[1].Size) == 0x00 {
		BoldRed.Println("\n[!] ERROR: File has empty import table.")
		os.Exit(1)
	}
	progress()

	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+peid.fileName+"|tr -d \""+peid.fileName+"\""+"|tr -d \"\n\"")).Output()
	ParseError(wcErr,"\n[!] ERROR: While getting the file size",string(wc))

	peid.fileSize = string(wc)
	progress()

	peid.Opt = Opt

	if peid.verbose == true {
		BoldYellow.Println("[*] File Size: " + peid.fileSize)
		BoldYellow.Printf("[*] Machine: %X\n", file.FileHeader.Machine)
		BoldYellow.Printf("[*] Magic: %X\n", Opt.Magic)
		BoldYellow.Printf("[*] Subsystem: %X\n", Opt.Subsystem)
		BoldYellow.Printf("[*] Image Base: %X\n", peid.imageBase)
		BoldYellow.Printf("[*] Size Of Image: %X\n", Opt.SizeOfImage)
		BoldYellow.Printf("[*] Export Table: %X\n", (Opt.DataDirectory[0].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Import Table: %X\n", (Opt.DataDirectory[1].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Base Relocation Table: %X\n", (Opt.DataDirectory[5].VirtualAddress + Opt.ImageBase))
		BoldYellow.Printf("[*] Import Address Table: %X\n", (Opt.DataDirectory[12].VirtualAddress + Opt.ImageBase))
	}
}

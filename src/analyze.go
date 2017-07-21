package main

import "debug/pe"
import "os/exec"

import "fmt"
import "os"

func analyze(file *pe.File) {
	//Do analysis on pe file...

	if file.FileHeader.Machine != 0x14C {
		boldRed.Println("\n[!] ERROR: File is not a 32 bit PE.")
		os.Exit(1)
	}
	progress()
	var OPT *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if OPT.Magic != 0x10B {
		boldRed.Println("\n[!] ERROR: File is not a valid PE.")
		os.Exit(1)
	}
	progress()
	peid.imageBase = OPT.ImageBase
	progress()
	peid.subsystem = OPT.Subsystem
	progress()

	if (OPT.DataDirectory[5].Size) != 0x00 {
		peid.aslr = true
		if peid.verbose == true {
			boldGreen.Println("[+] ASLR supported !")
			boldYellow.Println("[*] Using ASLR stub...")
		}
	} else {
		peid.aslr = false
		if peid.verbose == true {
			boldYellow.Println("[-] ASLR not supported :(")
			boldYellow.Println("[*] Using Non-ASLR stub...")
		}
	}
	progress()

	if (OPT.DataDirectory[11].Size) != 0x00 {
		boldRed.Println("\n[!] ERROR: File has bounded imports.")
		os.Exit(1)
	}
	progress()
	if (OPT.DataDirectory[13].Size) != 0x00 {
		boldRed.Println("\n[!] ERROR: File has delayed imports.")
		os.Exit(1)
	}
	progress()

	if (OPT.DataDirectory[1].Size) == 0x00 {
		boldRed.Println("\n[!] ERROR: File has empty import table.")
		os.Exit(1)
	}
	progress()

	wc, wcErr := exec.Command("sh", "-c", string("wc -c "+peid.fileName+"|tr -d \""+peid.fileName+"\""+"|tr -d \"\n\"")).Output()
	if wcErr != nil {
		boldRed.Println("\n[!] ERROR: While getting the file size")
		boldRed.Println(string(wc))
		fmt.Println(wcErr)
		clean()
		os.Exit(1)
	}

	peid.fileSize = string(wc)
	progress()

	peid.OPT = OPT

	if peid.verbose == true {
		boldYellow.Println("[*] File Size: " + peid.fileSize)
		boldYellow.Printf("[*] Machine: %X\n", file.FileHeader.Machine)
		boldYellow.Printf("[*] Magic: %X\n", OPT.Magic)
		boldYellow.Printf("[*] Subsystem: %X\n", OPT.Subsystem)
		boldYellow.Printf("[*] Image Base: %X\n", peid.imageBase)
		boldYellow.Printf("[*] Size Of Image: %X\n", OPT.SizeOfImage)
		boldYellow.Printf("[*] Export Table: %X\n", (OPT.DataDirectory[0].VirtualAddress + OPT.ImageBase))
		boldYellow.Printf("[*] Import Table: %X\n", (OPT.DataDirectory[1].VirtualAddress + OPT.ImageBase))
		boldYellow.Printf("[*] Base Relocation Table: %X\n", (OPT.DataDirectory[5].VirtualAddress + OPT.ImageBase))
		boldYellow.Printf("[*] Import Address Table: %X\n", (OPT.DataDirectory[12].VirtualAddress + OPT.ImageBase))
	}
}

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
	var Opt *pe.OptionalHeader32 = file.OptionalHeader.(*pe.OptionalHeader32)
	// PE32 = 0x10B
	if Opt.Magic != 0x10B {
		boldRed.Println("\n[!] ERROR: File is not a valid PE.")
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

	if (Opt.DataDirectory[11].Size) != 0x00 {
		boldRed.Println("\n[!] ERROR: File has bounded imports.")
		os.Exit(1)
	}
	progress()
	if (Opt.DataDirectory[13].Size) != 0x00 {
		boldRed.Println("\n[!] ERROR: File has delayed imports.")
		os.Exit(1)
	}
	progress()

	if (Opt.DataDirectory[1].Size) == 0x00 {
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

	peid.Opt = Opt

	if peid.verbose == true {
		boldYellow.Println("[*] File Size: " + peid.fileSize)
		boldYellow.Printf("[*] Machine: %X\n", file.FileHeader.Machine)
		boldYellow.Printf("[*] Magic: %X\n", Opt.Magic)
		boldYellow.Printf("[*] Subsystem: %X\n", Opt.Subsystem)
		boldYellow.Printf("[*] Image Base: %X\n", peid.imageBase)
		boldYellow.Printf("[*] Size Of Image: %X\n", Opt.SizeOfImage)
		boldYellow.Printf("[*] Export Table: %X\n", (Opt.DataDirectory[0].VirtualAddress + Opt.ImageBase))
		boldYellow.Printf("[*] Import Table: %X\n", (Opt.DataDirectory[1].VirtualAddress + Opt.ImageBase))
		boldYellow.Printf("[*] Base Relocation Table: %X\n", (Opt.DataDirectory[5].VirtualAddress + Opt.ImageBase))
		boldYellow.Printf("[*] Import Address Table: %X\n", (Opt.DataDirectory[12].VirtualAddress + Opt.ImageBase))
	}
}

/*
func MapPE(File string) {

	var Needle int = 0
	var SectionHeader *pe.SectionHeader32 = peid.Opt

	RawFile, readErr := ioutil.ReadFile(File)
	if readErr != nil {
		boldRed.Println("[!] ERROR: While reading the file")
		fmt.Print(readErr)
		os.Exit(1)
	}
	exec.Command("sh", "-c", "rm Mem.map").Run()
	Map, _ := os.Create("Mem.map")

	Map.Write(RawFile[0:peid.Opt.SizeOfHeaders])
	Needle += peid.Opt.SizeOfHeaders

}
*/
/*
func parseIAT(MemMap []byte) {

	var ImportTable uint32 = (peid.Opt.DataDirectory[1].VirtualAddress)
	var ImportAddressTable uint32 = (peid.Opt.DataDirectory[12].VirtualAddress)

	var IMAGE_IMPORT_DESCRIPTOR *pe.ImportDirectory

	IMAGE_IMPORT_DESCRIPTOR = unsafe.Pointer(&MemMap[ImportTable])




}
*/

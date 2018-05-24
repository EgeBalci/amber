package main

import "encoding/hex"
import "io/ioutil"
import "debug/pe"
import "errors"
import "bytes"

func CreateFileMapping(file string) (bytes.Buffer) {

	verbose("Mapping PE File...", "*")
	// Open the file as a *pe.File
	File, err := pe.Open(file)
	ParseError(err,"While opening file for mapping")
	progress()
	// Open the file as a byte array
	RawFile, err2 := ioutil.ReadFile(file)
	ParseError(err2,"While reading file content")
	progress()

	// if File.Machine == 0x8664 {
	// 	_opt := (File.OptionalHeader.(*pe.OptionalHeader64))
	// 	target.opt.Magic = _opt.Magic
	// 	target.opt.Subsystem = _opt.Subsystem
	// 	target.opt.CheckSum = _opt.CheckSum
	// 	target.opt.ImageBase = _opt.ImageBase
	// 	target.opt.AddressOfEntryPoint = _opt.AddressOfEntryPoint
	// 	target.opt.SizeOfImage =  _opt.SizeOfImage
	// 	target.opt.SizeOfHeaders = _opt.SizeOfHeaders
	// 	for i:=0; i<16; i++ {
	// 		target.opt.DataDirectory[i].VirtualAddress = _opt.DataDirectory[i].VirtualAddress
	// 		target.opt.DataDirectory[i].Size = _opt.DataDirectory[i].Size
	// 	}
	// }else{
	// 	_opt := File.OptionalHeader.((*pe.OptionalHeader32))
	// 	target.opt.Magic = _opt.Magic
	// 	target.opt.Subsystem = _opt.Subsystem
	// 	target.opt.CheckSum = _opt.CheckSum
	// 	target.opt.ImageBase = uint64(_opt.ImageBase)
	// 	target.opt.AddressOfEntryPoint = _opt.AddressOfEntryPoint
	// 	target.opt.SizeOfImage =  _opt.SizeOfImage
	// 	target.opt.SizeOfHeaders = _opt.SizeOfHeaders
	// 	for i:=0; i<16; i++ {
	// 		target.opt.DataDirectory[i].VirtualAddress = _opt.DataDirectory[i].VirtualAddress
	// 		target.opt.DataDirectory[i].Size = _opt.DataDirectory[i].Size
	// 	}
	// }

	// Check if the PE file is 64 bit (Will be removed)
	if File.Machine == 0x8664 {
		err := errors.New("64 bit files not supported.")
		ParseError(err,"Amber currently does not support 64 bit PE files.")
	}

	var offset uint64 = target.opt.ImageBase
	Map := bytes.Buffer{}
	Map.Write(RawFile[0:target.opt.SizeOfHeaders])
	offset += uint64(target.opt.SizeOfHeaders)
	progress()

	for i := 0; i < len(File.Sections); i++ {
		// Append null bytes if there is a gap between sections or PE header
		for {
			if offset < (uint64(File.Sections[i].VirtualAddress)+target.opt.ImageBase) {
				Map.WriteString(string(0x00))
				offset += 1
			} else {
				break
			}
		}
		// Map the section
		SectionData, err := File.Sections[i].Data()
		if err != nil {
			err := errors.New("Cannot read section data.")
			ParseError(err,"While reading the file section data.")
		}
		Map.Write(SectionData)
		offset += uint64(File.Sections[i].Size)
		// Append null bytes until reaching the end of the virtual address of the section
		for {
			if offset < (uint64(File.Sections[i].VirtualAddress)+uint64(File.Sections[i].VirtualSize)+target.opt.ImageBase) {
				Map.WriteString(string(0x00))
				offset += 1
			} else {
				break
			}
		}

	}
	progress()
	for {
		if (offset-target.opt.ImageBase) < uint64(target.opt.SizeOfImage) {
			Map.WriteString(string(0x00))
			offset += 1
		} else {
			break
		}
	}
	progress()
	
	// Perform integrity checks...
	verbose("\n[#] Performing integrity checks  on file mapping...\n|", "Y")
	if int(target.opt.SizeOfImage) != Map.Len() {
		if !target.IgnoreIntegrity {
			err := errors.New("Integrity check failed (Mapping size does not match the size of image header)\nTry '-ignore-integrity' parameter.")
			ParseError(err,"Integrity check failed (Mapping size does not match the size of image header)")
		}
	}
	verbose("[Image Size]------------> OK", "Y")

	for i := 0; i < len(File.Sections); i++ {
		for j := 0; j < int(File.Sections[i].Size/10); j++ {
			Buffer := Map.Bytes()
			if RawFile[int(int(File.Sections[i].Offset)+j)] != Buffer[int(int(File.Sections[i].VirtualAddress)+j)] {
				if !target.IgnoreIntegrity {
					err := errors.New("Integrity check failed (Broken section alignment)\nTry '-ignore-integrity' parameter.")
					ParseError(err,"Integrity check failed (Broken section alignment)")
				}
			}
		}
	}
	verbose("[Section Alignment]-----> OK\n", "Y")
	// Add data directory intervals check !
	progress()

	return Map
}


func scrape(Map []byte) ([]byte){

	verbose("Scraping PE headers...","*")

	var scraped string

	// if string(Map[:2]) == "MZ" {
	// 	scraped += hex.Dump(Map[:2])
	// 	Map[0] = 0x00
	// 	Map[1] = 0x00
	// }

	// for i:=0; i<0x1000; i++ {
	// 	if string(Map[i:i+2]) == "PE" {
	// 		scraped += hex.Dump(Map[i:i+2])
	// 		Map[i] = 0x00
	// 		Map[i+1] = 0x00
	// 	}
	// }

	for i:=0; i<0x1000; i++ {
		if string(Map[i:i+39]) == "This program cannot be run in DOS mode." {
			scraped += hex.Dump(Map[i:i+39])
			for j:=0; j<39; j++ {
				Map[i+j] = 0x00
			}
		}
	}

	for i:=66; i<0x1000; i++ {
		if Map[i] == 0x2e && Map[i+1] < 0x7e && Map[i+1] > 0x21 {
			scraped += hex.Dump(Map[i:i+7])
			for j:=0; j<7; j++{
				Map[i+j] = 0x00
			}
		}
	}

	verbose(scraped,"")
	verbose("Done scraping headers.","+")

	return Map
}
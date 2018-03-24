package main

import "encoding/hex"
import "io/ioutil"
import "debug/pe"
import "errors"
import "bytes"

func CreateFileMapping(file string) (bytes.Buffer) {

	verbose("Mapping PE file...", "*")
	// Open the file as a *pe.File
	File, err := pe.Open(file)
	ParseError(err,"While opening file for mapping")
	progress()
	// Open the file as a byte array
	RawFile, err2 := ioutil.ReadFile(file)
	ParseError(err2,"While reading file content")
	progress()
	OptionalHeader := File.OptionalHeader.(*pe.OptionalHeader32)
	// Check if the PE file is 64 bit
	if File.Machine == 0x8664 {
		err := errors.New("64 bit files not supported.")
		ParseError(err,"Amber currently does not support 64 bit PE files.")
	}
	var Offset uint32 = OptionalHeader.ImageBase
	Map := bytes.Buffer{}
	// Map the PE headers
	Map.Write(RawFile[0:int(OptionalHeader.SizeOfHeaders)])
	Offset += OptionalHeader.SizeOfHeaders
	progress()

	for i := 0; i < len(File.Sections); i++ {
		// Append null bytes if there is a gap between sections or PE header
		for {
			if Offset < (File.Sections[i].VirtualAddress + OptionalHeader.ImageBase) {
				Map.WriteString(string(0x00))
				Offset += 1
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
		Offset += File.Sections[i].Size
		// Append null bytes until reaching the end of the virtual address of the section
		for {
			if Offset < (File.Sections[i].VirtualAddress + File.Sections[i].VirtualSize + OptionalHeader.ImageBase) {
				Map.WriteString(string(0x00))
				Offset += 1
			} else {
				break
			}
		}

	}
	progress()
	for {
		if (Offset - OptionalHeader.ImageBase) < OptionalHeader.SizeOfImage {
			Map.WriteString(string(0x00))
			Offset += 1
		} else {
			break
		}
	}
	progress()
	// Perform integrity checks...
	verbose("\n[#] Performing integrity checks  on file mapping...\n|", "Y")
	if int(OptionalHeader.SizeOfImage) != Map.Len() {
		if !target.IgnoreMappingSize {
			err := errors.New("Integrity check failed (Mapping size does not match the size of image header)\nTry '--ignore-mapping-size' parameter.")
			ParseError(err,"Integrity check failed (Mapping size does not match the size of image header)")
		}
	}
	verbose("[Image Size]------------> OK", "Y")
	/*

		if Offset != ((File.Sections[len(File.Sections)-1].SectionHeader.VirtualAddress)+(File.Sections[len(File.Sections)-1].SectionHeader.VirtualSize)){
			err := errors.New("Integrity check failed (Offset does not match the final address)")
			return bytes.Buffer{},err
		}

	*/
	for i := 0; i < len(File.Sections); i++ {
		for j := 0; j < int(File.Sections[i].Size/10); j++ {
			Buffer := Map.Bytes()
			if RawFile[int(int(File.Sections[i].Offset)+j)] != Buffer[int(int(File.Sections[i].VirtualAddress)+j)] {
				if !target.IgnoreSectionAlignment {
					err := errors.New("Integrity check failed (Broken section alignment)\nTry '--ignore-section-alignment' parameter.")
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

	if string(Map[:2]) == "MZ" {
		verbose(hex.Dump(Map[:2]),"")
		Map[0] = byte(0x00)
		Map[1] = byte(0x00)
	}

	if string(Map[64:66]) == "PE" {
		verbose(hex.Dump(Map[64:66]),"")
		Map[64] = byte(0x00)
		Map[65] = byte(0x00)
	}
	
	if string(Map[78:117]) == "This program cannot be run in DOS mode." {
		verbose(hex.Dump(Map[78:117]),"")
		for i:=0; i<40; i++ {
			Map[78+i] = byte(0x00)
		}
	}

	if string(Map[128:130]) == "PE" {
		verbose(hex.Dump(Map[128:130]),"")
		Map[128] = byte(0x00)
		Map[129] = byte(0x00)
	}

	for i:=66; i<0x1000; i++{
		if Map[i] == 0x2e && Map[i+1] < 0x7e && Map[i+1] > 0x21 {
			verbose(hex.Dump(Map[i:i+7]),"")
			for j:=0; j<7; j++{
				Map[i+j] = byte(0x00)
			}
		}
	}

	verbose("Done scraping headers.","+")

	return Map
}
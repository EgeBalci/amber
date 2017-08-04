package Amber

import "io/ioutil"
import "debug/pe"
import "errors"
import "bytes"


func CreateFileMapping(file string) (*bytes.Buffer,error){
	
	// Open the file as a *pe.File
	File, err := pe.Open(file)
	if err != nil {
		return nil,err
	}

	// Open the file as a byte array
	RawFile, err2 := ioutil.ReadFile(file)
	if err2 != nil {
		return nil,err2
	}

	OptionalHeader := File.OptionalHeader.(*pe.OptionalHeader32)

	// Check if the PE file is 64 bit
	if File.Machine == 0x8664 {
		OptionalHeader = File.OptionalHeader.(*pe.OptionalHeader64)
	}

	var Offset uint32 = OptionalHeader.ImageBase

	Map := bytes.NewBuffer()
	// Map the PE headers
	Map.Write(RawFile[0:int(OptionalHeader.SizeOfHeaders)])
	Offset += OptionalHeader.SizeOfHeaders

	for i := 0; i < len(File.Sections); i++ {
		// Append null bytes if there is a gap between sections or PE header
		for ;; {
			if Offset < (File.Sections[i].VirtualAddress+OptionalHeader.ImageBase) {
				Map.Write(0x00)
				Offset += 1
			}else{
				break
			}
		}
		// Map the section
		SectionData, err := File.Sections[i].Data()
		if err != nil {
			err := errors.New("ERROR: Cannot read section data")
			return nil,err
		}
		Map.Write(SectionData)
		Offset += File.Sections[i].Size
		// Append null bytes until reaching the end of the virtual address of the section
		for ;; {
			if Offset < (File.Sections[i].VirtualAddress+File.Sections[i].VirtualSize+OptionalHeader.ImageBase) {
				Map.Write(0x00)
				Offset += 1				
			}else{
				break
			}
		}

	}

	// Perform integrity checks...

	if OptionalHeader.SizeOfImage != len(Map) {
		err := errors.New("ERROR: Integrity check failed (Mapping size does not match the size of image header)")
		return nil,err	
	}


	if Offset != (File.Sections[len(File.Sections-1)].SectionHeader.VirtualAddress+File.Sections[len(File.Sections-1)].SectionHeader.VirtualSize){
		err := errors.New("ERROR: Integrity check failed (Offset does not match the final address)")
		return nil,err
	}


	for i := 0; i < len(File.Sections); i++ {
		for j := 0; j < len(File.Sections[i].Size/10); j++ {

			if RawFile[int(File.Sections[i].Offset+j)] != Map[File.Sections[i].VirtualAddress+j]{
				err := errors.New("ERROR: Integrity check failed (Broken section alignment)")
				return nil,err
			}
		}
	}





	return Map,nil

}
package main

import "io/ioutil"
import "debug/pe"
import "errors"
import "bytes"


func CreateFileMapping(file string) (bytes.Buffer,error){

	verbose("[*] Mapping PE file...",BoldYellow)

	// Open the file as a *pe.File
	File, err := pe.Open(file)
	if err != nil {
		return bytes.Buffer{},err
	}

	progress()

	// Open the file as a byte array
	RawFile, err2 := ioutil.ReadFile(file)
	if err2 != nil {
		return bytes.Buffer{},err2
	}

	progress()

	OptionalHeader := File.OptionalHeader.(*pe.OptionalHeader32)

	// Check if the PE file is 64 bit
	if File.Machine == 0x8664 {
		err := errors.New("[!] ERROR: 64 bit files not supported.")
		return bytes.Buffer{},err
	}

	var Offset uint32 = OptionalHeader.ImageBase

	Map := bytes.Buffer{}
	// Map the PE headers
	Map.Write(RawFile[0:int(OptionalHeader.SizeOfHeaders)])
	Offset += OptionalHeader.SizeOfHeaders

	progress()


	for i := 0; i < len(File.Sections); i++ {
		// Append null bytes if there is a gap between sections or PE header
		for ;; {
			if Offset < (File.Sections[i].VirtualAddress+OptionalHeader.ImageBase) {
				Map.WriteString(string(0x00))
				Offset += 1
			}else{
				break
			}
		}
		// Map the section
		SectionData, err := File.Sections[i].Data()
		if err != nil {
			err := errors.New("[!] ERROR: Cannot read section data")
			return bytes.Buffer{},err
		}
		Map.Write(SectionData)
		Offset += File.Sections[i].Size
		// Append null bytes until reaching the end of the virtual address of the section
		for ;; {
			if Offset < (File.Sections[i].VirtualAddress+File.Sections[i].VirtualSize+OptionalHeader.ImageBase) {
				Map.WriteString(string(0x00))
				Offset += 1
			}else{
				break
			}
		}

	}

	progress()

	for ;; {
		if (Offset-OptionalHeader.ImageBase) < OptionalHeader.SizeOfImage {
			Map.WriteString(string(0x00))
			Offset += 1
		}else{
			break
		}		
	}
	
	progress()

	// Perform integrity checks...

	verbose("\n[#] Performing integrity checks  on file mapping...\n|",BoldYellow)

	if int(OptionalHeader.SizeOfImage) != Map.Len() {
		err := errors.New("[!] ERROR: Integrity check failed (Mapping size does not match the size of image header)")
		return bytes.Buffer{},err
	}

	verbose("[Image Size]------------> OK",BoldYellow)
/*

	if Offset != ((File.Sections[len(File.Sections)-1].SectionHeader.VirtualAddress)+(File.Sections[len(File.Sections)-1].SectionHeader.VirtualSize)){
		err := errors.New("[!] ERROR: Integrity check failed (Offset does not match the final address)")
		return bytes.Buffer{},err
	}

*/
	for i := 0; i < len(File.Sections); i++ {
		for j := 0; j < int(File.Sections[i].Size/10); j++ {

			Buffer := Map.Bytes()

			if RawFile[int(int(File.Sections[i].Offset)+j)] != Buffer[int(int(File.Sections[i].VirtualAddress)+j)]{
				err := errors.New("[!] ERROR: Integrity check failed (Broken section alignment)")
				return bytes.Buffer{},err
			}
		}
	}

	verbose("[Section Alignment]-----> OK\n",BoldYellow)

	// Add data directory intervals check !

	progress()



	return Map,nil

}

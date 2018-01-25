package main

import "os/exec"
import "os"

func assemble() {

	verbose("Assembling reflective payload...","*")

	// Create a file mapping image (6 steps)
	Map, MapErr:= CreateFileMapping(peid.FileName)
	ParseError(MapErr,"While creating file mapping","")
	MapFile, MapFileErr := os.Create("Mem.map")
	ParseError(MapFileErr,"While getting the file size","")

	MapFile.Write(Map.Bytes())
	MapFile.Close()

	progress()

	if peid.aslr == false {
		if peid.iat == true {
			move("/usr/share/Amber/Mem.map","/usr/share/Amber/core/Fixed/iat/Mem.map")
		}else{
			move("/usr/share/Amber/Mem.map","/usr/share/Amber/core/Fixed/Mem.map")
		}
		progress()
		Cdir("/usr/share/Amber/core/Fixed")
		if peid.iat == true {
			Cdir("/usr/share/Amber/core/Fixed/iat")		
		}
		progress()
		nasm, Err := exec.Command("nasm","-f","bin","Stub.asm","-o","/usr/share/Amber/Payload").Output()
		ParseError(Err,"While assembling payload :(",string(nasm))
		progress()
	} else {
		if peid.iat == true {
			move("/usr/share/Amber/Mem.map","/usr/share/Amber/core/ASLR/iat/Mem.map")
		}else{
			move("/usr/share/Amber/Mem.map","/usr/share/Amber/core/ASLR/Mem.map")
		}
		progress()
		Cdir("/usr/share/Amber/core/ASLR")		
		if peid.iat == true {
			Cdir("/usr/share/Amber/core/ASLR/iat/")
		}
		progress()
		nasm, Err := exec.Command("nasm","-f","bin","Stub.asm","-o","/usr/share/Amber/Payload").Output()
		ParseError(Err,"While assembling payload :(",string(nasm))
		progress()
	}

	Cdir("/usr/share/Amber")

	verbose("Assebly completed.", "*")
}

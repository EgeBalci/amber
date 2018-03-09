package main

import "os/exec"
import "os"

func assemble() {

	verbose("Assembling reflective payload...", "*")
	// Create a file mapping image (6 steps)
	Map := CreateFileMapping(target.FileName)
	MapFile, MapFileErr := os.Create("Mem.map")
	target.clean = true
	ParseError(MapFileErr, "While getting the file size")

	MapFile.Write(Map.Bytes())
	MapFile.Close()
	progress()

	if target.aslr == false {
		if target.iat == true {
			move("/usr/share/Amber/Mem.map", "/usr/share/Amber/core/Fixed/iat/Mem.map")
		} else {
			move("/usr/share/Amber/Mem.map", "/usr/share/Amber/core/Fixed/Mem.map")
		}
		progress()
		Cdir("/usr/share/Amber/core/Fixed")
		if target.iat == true {
			Cdir("/usr/share/Amber/core/Fixed/iat")
		}
		progress()
		Err := exec.Command("nasm", "-f", "bin", "stub.asm", "-o", "/usr/share/Amber/Payload").Run()
		ParseError(Err, "While assembling payload :(")
		progress()
	} else {
		if target.iat == true {
			move("/usr/share/Amber/Mem.map", "/usr/share/Amber/core/ASLR/iat/Mem.map")
		} else {
			move("/usr/share/Amber/Mem.map", "/usr/share/Amber/core/ASLR/Mem.map")
		}
		progress()
		Cdir("/usr/share/Amber/core/ASLR")
		if target.iat == true {
			Cdir("/usr/share/Amber/core/ASLR/iat/")
		}
		progress()
		Err := exec.Command("nasm", "-f", "bin", "stub.asm", "-o", "/usr/share/Amber/Payload").Run()
		ParseError(Err, "While assembling payload :(")
		progress()
	}

	Cdir("/usr/share/Amber")
	verbose("Assebly completed.", "*")
}

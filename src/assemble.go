package main

import "os/exec"
import "os"

func assemble() {

	verbose("[*] Assembling reflective payload...",BoldYellow)

	// Create a file mapping image (6 steps)
	Map, MapErr:= CreateFileMapping(peid.fileName)
	ParseError(MapErr,"\n[!] ERROR: While creating file mapping","")
	MapFile, MapFileErr := os.Create("Mem.map")
	ParseError(MapFileErr,"\n[!] ERROR: While getting the file size","")

	MapFile.Write(Map.Bytes())
	MapFile.Close()

	progress()

	if peid.aslr == false {
		moveMapCommand := "mv Mem.map core/Fixed/"
		if peid.iat == true {
			moveMapCommand += "iat/"
		}
		moveMap, moveMapErr := exec.Command("sh", "-c", moveMapCommand).Output()
		ParseError(moveMapErr,"\n[!] ERROR: While moving the file map",string(moveMap))
		progress()
		nasmCommand := "cd core/Fixed/"
		if peid.iat == true {
			nasmCommand += "iat/"
		}
		nasmCommand += " && nasm -f bin Stub.asm -o Payload"
		nasm, Err := exec.Command("sh", "-c", nasmCommand).Output()
		ParseError(Err,"\n[!] ERROR: While assembling payload :(",string(nasm))

		progress()
		movePayloadCommand := "mv core/Fixed/Payload ./"
		if peid.iat == true {
			movePayloadCommand = "mv core/Fixed/iat/Payload ./"
		}
		movePayload, movePayErr := exec.Command("sh", "-c", movePayloadCommand).Output()
		ParseError(movePayErr,"\n[!] ERROR: While moving the payload",string(movePayload))

		progress()
	} else {
		moveMapCommand := "mv Mem.map core/ASLR/"
		if peid.iat == true {
			moveMapCommand += "iat/"
		}
		moveMap, moveMapErr := exec.Command("sh", "-c", moveMapCommand).Output()
		ParseError(moveMapErr,"\n[!] ERROR: While moving the file map",string(moveMap))
		progress()
		nasmCommand := "cd core/ASLR/"
		if peid.iat == true {
			nasmCommand += "iat/"
		}
		nasmCommand += " && nasm -f bin Stub.asm -o Payload"
		nasm, Err := exec.Command("sh", "-c", nasmCommand).Output()
		ParseError(Err,"\n[!] ERROR: While assembling payload :(",string(nasm))

		progress()
		movePayloadCommand := "mv core/ASLR/Payload ./"
		if peid.iat == true {
			movePayloadCommand = "mv core/ASLR/iat/Payload ./"
		}
		movePayload, movePayErr := exec.Command("sh", "-c", movePayloadCommand).Output()
		ParseError(movePayErr,"\n[!] ERROR: While moving the payload",string(movePayload))
		progress()
	}
	verbose("[*] Assebly completed.", yellow)
}

package main

import "os/exec"
import "os"

func assemble() {

	// Create a file mapping image (5 steps)
	Map, MapErr:= CreateFileMapping(peid.fileName)
	MapFile, MapFileErr := os.Create("Mem.map")
	if MapFileErr != nil || MapErr != nil{
		BoldRed.Println("\n[!] ERROR: While creating file mapping")
		BoldRed.Println(MapFileErr)
		BoldRed.Println(MapErr)
		clean()
		os.Exit(1)
	}

	MapFile.Write(Map.Bytes())
	MapFile.Close()


	progress()
	if peid.aslr == false {
		moveMapCommand := "mv Mem.map core/NonASLR/"
		if peid.iat == true {
			moveMapCommand += "iat/"
		}
		moveMap, moveMapErr := exec.Command("sh", "-c", moveMapCommand).Output()
		if moveMapErr != nil {
			BoldRed.Println("\n[!] ERROR: While moving the file map")
			BoldRed.Println(string(moveMap))
			clean()
			os.Exit(1)
		}
		progress()
		nasmCommand := "cd core/NonASLR/"
		if peid.iat == true {
			nasmCommand += "iat/"
		}
		nasmCommand += " && nasm -f bin Stub.asm -o Payload"
		nasm, Err := exec.Command("sh", "-c", nasmCommand).Output()
		if Err != nil {
			BoldRed.Println("\n[!] ERROR: While assembling payload :(")
			BoldRed.Println(string(nasm))
			BoldRed.Println(Err)
			clean()
			os.Exit(1)
		}
		progress()
		movePayloadCommand := "mv core/NonASLR/Payload ./"
		if peid.iat == true {
			movePayloadCommand = "mv core/NonASLR/iat/Payload ./"
		}
		movePayload, movePayErr := exec.Command("sh", "-c", movePayloadCommand).Output()
		if movePayErr != nil {
			BoldRed.Println("\n[!] ERROR: While moving the payload")
			BoldRed.Println(string(movePayload))
			BoldRed.Println(Err)
			clean()
			os.Exit(1)
		}
		progress()
	} else {
		moveMapCommand := "mv Mem.map core/ASLR/"
		if peid.iat == true {
			moveMapCommand += "iat/"
		}
		moveMap, moveMapErr := exec.Command("sh", "-c", moveMapCommand).Output()
		if moveMapErr != nil {
			BoldRed.Println("\n[!] ERROR: While moving the file map")
			BoldRed.Println(string(moveMap))
			clean()
			os.Exit(1)
		}
		progress()
		nasmCommand := "cd core/ASLR/"
		if peid.iat == true {
			nasmCommand += "iat/"
		}
		nasmCommand += " && nasm -f bin Stub.asm -o Payload"
		nasm, Err := exec.Command("sh", "-c", nasmCommand).Output()
		if Err != nil {
			BoldRed.Println("\n[!] ERROR: While assembling payload :(")
			BoldRed.Println(string(nasm))
			BoldRed.Println(Err)
			clean()
			os.Exit(1)
		}
		progress()
		movePayloadCommand := "mv core/ASLR/Payload ./"
		if peid.iat == true {
			movePayloadCommand = "mv core/ASLR/iat/Payload ./"
		}
		movePayload, movePayErr := exec.Command("sh", "-c", movePayloadCommand).Output()
		if movePayErr != nil {
			BoldRed.Println("\n[!] ERROR: While moving the payload")
			BoldRed.Println(string(movePayload))
			BoldRed.Println(Err)
			clean()
			os.Exit(1)
		}
		progress()
	}
	verbose("[*] Assebly completed.", yellow)
}

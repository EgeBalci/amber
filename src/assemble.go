package main

import "os/exec"
import "strings"
import "os"

func assemble() {

	MapPE, _ := exec.Command("sh", "-c", string("wine MapPE.exe "+peid.fileName)).Output()
	if strings.Contains(string(MapPE), "[!]") {
		BoldRed.Println("\n[!] ERROR: While mapping pe file :(")
		BoldRed.Println(string(MapPE))
		clean()
		os.Exit(1)
	}

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
	_MapPE := strings.Split(string(MapPE), "github.com/egebalci/mappe")
	verbose(string(_MapPE[1]), white)
}

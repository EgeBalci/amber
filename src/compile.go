package main

import "os/exec"

//import "os"

func compile() {

	move("Payload.rc4", "Payload")
	xxd_err := exec.Command("sh", "-c", "xxd -i Payload > stub/payload.h").Run()
	if xxd_err != nil {
		ParseError(xxd_err, "While extracting payload hex stream.", "")
	}
	progress()

	_xxd_err := exec.Command("sh", "-c", "xxd -i Payload.key > stub/key.h").Run()
	if _xxd_err != nil {
		ParseError(_xxd_err, "While extracting key hex stream.", "")
	}
	progress()

	var compileCommand string = "i686-w64-mingw32-g++-win32 -c stub/stub.cpp"
	if PACKET_MANAGER == "pacman" {
		compileCommand = "i686-w64-mingw32-g++ -c stub/stub.cpp"
	}

	mingwObj, mingwObjErr := exec.Command("sh", "-c", compileCommand).Output()
	ParseError(mingwObjErr, "While compiling the object file.", string(mingwObj))

	progress()

	compileCommand = "i686-w64-mingw32-g++-win32 stub.o "
	if PACKET_MANAGER == "pacman" {
		compileCommand = "i686-w64-mingw32-g++ stub.o "
	}

	if peid.resource == true {
		compileCommand += "stub/Resource.o "
	}
	if peid.subsystem == 2 { // GUI
		compileCommand += string("-mwindows -o " + peid.FileName)
	} else {
		compileCommand += string("-o " + peid.FileName)
	}
	progress()

	verbose("Compiling to EXE...", "*")
	mingw, mingwErr := exec.Command("sh", "-c", compileCommand).Output()
	ParseError(mingwErr, "While compiling to exe. (This might caused by a permission issue)", string(mingw))

	progress()

	strip, stripErr := exec.Command("sh", "-c", string("strip "+peid.FileName)).Output()
	ParseError(stripErr, "While striping the exe.", string(strip))

	progress()
}

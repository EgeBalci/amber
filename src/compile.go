package main

import "os/exec"

func compile() {

	move("Payload.rc4", "Payload")
	xxd_err := exec.Command("sh", "-c", "xxd -i Payload > stub/payload.h").Run()
	ParseError(xxd_err, "While extracting payload hex stream.")
	progress()

	_xxd_err := exec.Command("sh", "-c", "xxd -i Payload.key > stub/key.h").Run()
	ParseError(_xxd_err, "While extracting key hex stream.")
	progress()

	var compileCommand string = "i686-w64-mingw32-g++-win32 -c stub/stub.cpp"
	if PACKET_MANAGER == "pacman" {
		compileCommand = "i686-w64-mingw32-g++ -c stub/stub.cpp"
	}

	mingwObjErr := exec.Command("sh", "-c", compileCommand).Run()
	ParseError(mingwObjErr, "While compiling the object file.",)
	progress()

	compileCommand = "i686-w64-mingw32-g++-win32 stub.o "
	if PACKET_MANAGER == "pacman" {
		compileCommand = "i686-w64-mingw32-g++ stub.o "
	}

	if target.resource == false {
		compileCommand += "stub/Resource.o "
		verbose("Adding resource data...", "*")
	}
	if target.subsystem == 2 { // GUI
		compileCommand += string("-mwindows -o " + target.FileName)
	} else {
		compileCommand += string("-o " + target.FileName)
	}
	progress()

	verbose("Compiling to EXE...", "*")
	mingwErr := exec.Command("sh", "-c", compileCommand).Run()
	ParseError(mingwErr, "While compiling to exe. (This might caused by a permission issue)")
	progress()

	stripErr := exec.Command("sh", "-c", string("strip "+target.FileName)).Run()
	ParseError(stripErr, "While striping the exe.")
	progress()
}

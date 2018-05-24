package main

import "strconv"
import "os/exec"

func compile() {

	move("Payload.rc4", "Payload")
	xxd_err := exec.Command("sh", "-c", "xxd -i Payload > stub/payload.h").Run()
	ParseError(xxd_err, "While extracting payload hex stream.")
	progress()

	xxd_err = exec.Command("sh", "-c", "xxd -i Payload.key > stub/key.h").Run()
	ParseError(xxd_err, "While extracting key hex stream.")
	progress()

	var CompileCommand string = "i686-w64-mingw32-g++-win32 -c stub/stub.cpp"
	if PACKET_MANAGER == "pacman" {
		CompileCommand = "i686-w64-mingw32-g++ -c stub/stub.cpp"
	}

	MingwObjErr := exec.Command("sh", "-c", CompileCommand).Run()
	ParseError(MingwObjErr, "While compiling the object file.",)
	progress()

	ImageBase := strconv.FormatInt(int64(target.opt.ImageBase), 16)
	CompileCommand = "i686-w64-mingw32-g++-win32 -Wl,--image-base,0x"+ImageBase+" stub.o "
	if PACKET_MANAGER == "pacman" {
		CompileCommand = "i686-w64-mingw32-g++ -Wl,--image-base,0x"+ImageBase+" stub.o "
	}

	if target.dll {
		CompileCommand += "-shared -o "+target.FileName
	}else{
		if target.resource == false {
			CompileCommand += "stub/Resource.o "
			verbose("Adding resource data...", "*")
		}
		if target.opt.Subsystem == 2 { // GUI
			CompileCommand += "-mwindows -o "+target.FileName
		}else{
			CompileCommand += "-o "+target.FileName
		}
	}
	//verbose(CompileCommand, "*")
	progress()
	verbose("Compiling to EXE...", "*")
	//verbose(CompileCommand,"*")
	MingwErr := exec.Command("sh", "-c", CompileCommand).Run()
	ParseError(MingwErr, "While compiling to exe. (This might caused by a permission issue)")
	progress()

	StripErr := exec.Command("sh", "-c", string("strip "+target.FileName)).Run()
	ParseError(StripErr, "While striping the exe.")
	progress()
}

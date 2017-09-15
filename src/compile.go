package main

import "os/exec"
import "os"

func compile() {
 
  verbose("[*] Ciphering payload...",BoldYellow)
  crypt() // 4 steps

  xxd := exec.Command("sh", "-c", "rm Payload && mv Payload.xor Payload && xxd -i Payload > stub/PAYLOAD.h")
  xxd.Stdout = os.Stdout
  xxd.Stderr = os.Stderr
  xxd.Run()
  progress()

  _xxd := exec.Command("sh", "-c", "xxd -i Payload.key > stub/KEY.h")
  _xxd.Stdout = os.Stdout
  _xxd.Stderr = os.Stderr
  _xxd.Run()
  progress()

  mingwObj, mingwObjErr := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 -c stub/stub.cpp").Output()
  ParseError(mingwObjErr,"\n[!] ERROR: While compiling the object file.",string(mingwObj))
  
  progress()

  var compileCommand string = "i686-w64-mingw32-g++-win32 stub.o "
  if peid.resource == true {
    compileCommand += "stub/Resource.o "
  }
  if peid.subsystem == 2 { // GUI
    compileCommand += string("-mwindows -o "+peid.fileName)
  }else{
    compileCommand += string("-o "+peid.fileName)
  }
  progress()

  verbose("[*] Compiling to EXE...",BoldYellow)
  mingw, mingwErr := exec.Command("sh", "-c", compileCommand).Output()
  ParseError(mingwErr,"\n[!] ERROR: While compiling to exe.",string(mingw))

  progress()

  strip, stripErr := exec.Command("sh", "-c", string("strip "+peid.fileName)).Output()
  ParseError(stripErr,"\n[!] ERROR: While striping the exe.",string(strip))

  progress()
}

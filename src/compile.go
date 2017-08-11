package main

import "os/exec"
import "fmt"
import "os"

func compile() {

 
  verbose("[*] Ciphering payload...",BoldYellow)
  crypt() // 4 steps

  xxd := exec.Command("sh", "-c", "rm Payload && mv Payload.xor Payload && xxd -i Payload > Stub/PAYLOAD.h")
  xxd.Stdout = os.Stdout
  xxd.Stderr = os.Stderr
  xxd.Run()
  progress()

  _xxd := exec.Command("sh", "-c", "xxd -i Payload.key > Stub/KEY.h")
  _xxd.Stdout = os.Stdout
  _xxd.Stderr = os.Stderr
  _xxd.Run()
  progress()

  mingwObj, mingwObjErr := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 -c Stub/Stub.cpp").Output()
  if mingwObjErr != nil {
    BoldRed.Println("\n[!] ERROR: While compiling the object file.")
    red.Println(string(mingwObj))
    fmt.Println(mingwObjErr)
    clean()
    os.Exit(1)
  }
  progress()

  var compileCommand string = "i686-w64-mingw32-g++-win32 Stub.o "
  if peid.resource == true {
    compileCommand += "Stub/Resource.o "
  }
  if peid.subsystem == 2 { // GUI
    compileCommand += string("-mwindows -o "+peid.fileName)
  }else{
    compileCommand += string("-o "+peid.fileName)
  }
  progress()

  mingw, mingwErr := exec.Command("sh", "-c", compileCommand).Output()
  if mingwErr != nil {
    BoldRed.Println("\n[!] ERROR: While compiling to exe.")
    red.Println(compileCommand)
    red.Println(string(mingw))
    fmt.Println(mingwErr)
    clean()
    os.Exit(1)
  }

  progress()
}

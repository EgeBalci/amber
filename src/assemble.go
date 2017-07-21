package main

import "os/exec"
import "strings"
import "os"

func assemble() {

  MapPE, _ := exec.Command("sh", "-c", string("wine MapPE.exe "+peid.fileName)).Output()
  if strings.Contains(string(MapPE), "[!]") {
    boldRed.Println("\n[!] ERROR: While mapping pe file :(")
    boldRed.Println(string(MapPE))
    clean()
    os.Exit(1)
  }
  progress()
  if peid.aslr == false {
    moveMap, moveMapErr := exec.Command("sh", "-c", "mv Mem.map Ophio/NonASLR/").Output()
    if moveMapErr != nil {
      boldRed.Println("\n[!] ERROR: While moving the file map")
      boldRed.Println(string(moveMap))
      clean()
      os.Exit(1)
    }
    progress()
    nasm, Err := exec.Command("sh", "-c", "cd Ophio/NonASLR/ && nasm -f bin Ophio.asm -o Payload").Output()
    if Err != nil {
      boldRed.Println("\n[!] ERROR: While assembling payload :(")
      boldRed.Println(string(nasm))
      boldRed.Println(Err)
      clean()
      os.Exit(1)
    }
    progress()
    movePayload, movePayErr := exec.Command("sh", "-c", "mv Ophio/NonASLR/Payload ./").Output()
    if movePayErr != nil {
      boldRed.Println("\n[!] ERROR: While moving the payload")
      boldRed.Println(string(movePayload))
      boldRed.Println(Err)
      clean()
      os.Exit(1)
    }
    progress()
  }else{
    moveMap, moveMapErr := exec.Command("sh", "-c", "mv Mem.map Ophio/ASLR/").Output()
    if moveMapErr != nil {
      boldRed.Println("\n[!] ERROR: While moving the file map")
      boldRed.Println(string(moveMap))
      clean()
      os.Exit(1)
    }
    progress()
    nasm, Err := exec.Command("sh", "-c", "cd Ophio/ASLR/ && nasm -f bin Ophio.asm -o Payload").Output()
    if Err != nil {
      boldRed.Println("\n[!] ERROR: While assembling payload :(")
      boldRed.Println(string(nasm))
      boldRed.Println(Err)
      clean()
      os.Exit(1)
    }
    progress()

    movePayload, movePayErr := exec.Command("sh", "-c", "mv Ophio/ASLR/Payload ./").Output()
    if movePayErr != nil {
      boldRed.Println("\n[!] ERROR: While moving the payload")
      boldRed.Println(string(movePayload))
      boldRed.Println(Err)
      clean()
      os.Exit(1)
    }
    progress()
  }
  _MapPE := strings.Split(string(MapPE), "github.com/egebalci/mappe")
  verbose(string(_MapPE[1]),white)
}
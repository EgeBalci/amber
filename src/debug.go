package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/cheggaaa/pb.v1"
)

func progress() {
	if peid.verbose == false {
		progressBar.Increment()
	}
}

func createBar() {

	var full int = 43

	if peid.verbose == false {
		if peid.staged == true {
			full -= 10
		}

		progressBar = pb.New(full)
		progressBar.SetWidth(80)
		progressBar.Start()
	}
}

func verbose(str string, col *color.Color) {

	if peid.verbose == true {
		col.Println(str)
	}

}

func checkRequired() {

	CheckMingw, mingwErr := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 --version").Output()
	if !strings.Contains(string(CheckMingw), "Copyright") {
		boldRed.Println("\n\n[!] ERROR: mingw is not installed.")
		red.Println(string(CheckMingw))
		red.Println(mingwErr)
		os.Exit(1)
	}
	progress()
	CheckNasm, _ := exec.Command("sh", "-c", "nasm -h").Output()
	if !strings.Contains(string(CheckNasm), "usage:") {
		boldRed.Println("\n\n[!] ERROR: nasm is not installed.")
		red.Println(string(CheckNasm))
		os.Exit(1)
	}
	progress()
	CheckStrip, _ := exec.Command("sh", "-c", "strip -V").Output()
	if !strings.Contains(string(CheckStrip), "Copyright") {
		boldRed.Println("\n\n[!] ERROR: strip is not installed.")
		red.Println(string(CheckStrip))
		os.Exit(1)
	}
	progress()
	CheckWine, _ := exec.Command("sh", "-c", "wine --help").Output()
	if !strings.Contains(string(CheckWine), "Usage:") {
		boldRed.Println("\n\n[!] ERROR: wine is not installed.")
		red.Println(string(CheckWine))
		os.Exit(1)
	}
	progress()
	CheckMapPE, _ := exec.Command("sh", "-c", "ls MapPE.exe").Output()
	if !strings.Contains(string(CheckMapPE), "MapPE.exe") {
		boldRed.Println("\n\n[!] ERROR: MapPE.exe is missing.")
		red.Println(string(CheckMapPE))
		red.Println(mingwErr)
		os.Exit(1)
	}
	progress()
	CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
	if !strings.Contains(string(CheckXXD), "Amber") {
		boldRed.Println("\n\n[!] ERROR: xxd is not installed.")
		red.Println(string(CheckMingw))
		os.Exit(1)
	}
	progress()
	CheckMultiLib, _ := exec.Command("sh", "-c", "apt-cache policy gcc-multilib").Output()
	if strings.Contains(string(CheckMultiLib), "(none)") {
		boldRed.Println("\n\n[!] ERROR: gcc-multilib is not installed.")
		red.Println(string(CheckMultiLib))
		os.Exit(1)
	}
	progress()
	CheckMultiLibPlus, _ := exec.Command("sh", "-c", "apt-cache policy g++-multilib").Output()
	if strings.Contains(string(CheckMultiLibPlus), "(none)") {
		boldRed.Println("\n\n[!] ERROR: g++-multilib is not installed.")
		red.Println(string(CheckMultiLibPlus))
		os.Exit(1)
	}
	progress()

}

func clean() {

	exec.Command("sh", "-c", "rm Ophio/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm Ophio/iat/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm Stub.o").Run()
	progress()
	exec.Command("sh", "-c", "rm Payload").Run()
	progress()
	exec.Command("sh", "-c", "rm Payload.key").Run()
	progress()
	exec.Command("sh", "-c", "echo   > Stub/PAYLOAD.h").Run()
	progress()
	exec.Command("sh", "-c", "echo   > Stub/KEY.h").Run()
	progress()

}

package main

import "io"
import "os"
import "fmt"
import "os/exec"
import "strings"
import "gopkg.in/cheggaaa/pb.v1"

func progress() {
	if target.verbose == false {
		progressBar.Increment()
	}
}

func CreateProgressBar() {

	var full int = 48

	if target.verbose == false {
		if target.staged == true {
			full -= 10
		}

		progressBar = pb.New(full)
		progressBar.SetWidth(80)
		progressBar.Start()
	}
}

func verbose(str string, status string) {

	if target.verbose {
		if status == "*" {
			BoldYellow.Print("[*] ")
			white.Println(str)
		} else if status == "+" {
			BoldGreen.Print("[+] ")
			white.Println(str)
		} else if status == "-" {
			BoldRed.Print("[-] ")
			white.Println(str)
		} else if status == "!" {
			BoldRed.Print("[!] ")
			white.Println(str)
		} else if status == "R" {
			BoldRed.Println(str)
		} else if status == "Y" {
			BoldYellow.Println(str)
		} else if status == "B" {
			BoldBlue.Println(str)
		} else if status == "" {
			fmt.Println(str)
		}
	}
}

func _verbose(str string, value uint64) {

	if target.verbose {
		BoldYellow.Print("[*] ")
		white.Printf(str+" 0x%X\n", value)
	}
}

func CheckRequirements() {

	verbose("Checking requirments...", "*")

	if PACKET_MANAGER == "pacman" {
		progress()
		CheckMingw, mingwErr := exec.Command("i686-w64-mingw32-g++", "--version").Output()
		if !strings.Contains(string(CheckMingw), "Copyright") {
			ParseError(mingwErr, "MingW is not installed correctly.")
		}
		progress()
		CheckNasm, _ := exec.Command("nasm", "-h").Output()
		if !strings.Contains(string(CheckNasm), "usage:") {
			ParseError(nil, "nasm is not installed correctly.")
		}
		progress()
		CheckStrip, _ := exec.Command("strip", "-V").Output()
		if !strings.Contains(string(CheckStrip), "Copyright") {
			ParseError(nil, "strip is not installed correctly.")
		}
		progress()
		CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
		if !strings.Contains(string(CheckXXD), "Amber") {
			ParseError(nil, "xxd is not installed correctly.")
		}
		progress()
		CheckMultiLib, _ := exec.Command("sh", "-c", "pacman -Qs gcc-multilib").Output()
		if strings.Contains(string(CheckMultiLib), "The GNU Compiler") {
			ParseError(nil, "gcc-multilib is not installed correctly.")
		}
		progress()
	} else {
		CheckMingw, mingwErr := exec.Command("i686-w64-mingw32-g++-win32", "--version").Output()
		if !strings.Contains(string(CheckMingw), "Copyright") {
			ParseError(mingwErr, "MingW is not installed correctly.")
		}
		progress()
		CheckNasm, _ := exec.Command("nasm", "-h").Output()
		if !strings.Contains(string(CheckNasm), "usage:") {
			ParseError(nil, "nasm is not installed correctly.")
		}
		progress()
		CheckStrip, _ := exec.Command("strip", "-V").Output()
		if !strings.Contains(string(CheckStrip), "Copyright") {
			ParseError(nil, "strip is not installed correctly.")
		}
		progress()
		CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
		if !strings.Contains(string(CheckXXD), "Amber") {
			ParseError(nil, "xxd is not installed correctly.")
		}
		progress()
		CheckMultiLib, _ := exec.Command("sh", "-c", "apt-cache policy gcc-multilib").Output()
		if strings.Contains(string(CheckMultiLib), "(none)") {
			ParseError(nil, "gcc-multilib is not installed correctly.")
		}
		progress()
		CheckMultiLibPlus, _ := exec.Command("sh", "-c", "apt-cache policy g++-multilib").Output()
		if strings.Contains(string(CheckMultiLibPlus), "(none)") {
			ParseError(nil, "g++-multilib is not installed correctly.")
		}
		progress()

	}
}

func ParseError(err error, ErrStatus string) {

	if err != nil {
		//progressBar.Finish()
		BoldRed.Println("\n\n\n[!] ERROR: " + ErrStatus)
		verbose("|\n|>","R")
		verbose(err.Error(),"")
		verbose("\n}\n","R")
		fmt.Println("\n")
		clean()
		os.Exit(1)

	}
}


func Cdir(dir string) {
	err := os.Chdir(dir)
	ParseError(err, "While changing directory.")
}

func move(old, new string) {
	err := os.Rename(old, new)
	if err != nil {
		ParseError(err, "While moving a file.")
	}
}

func _copy(old, new string) {
	from, err := os.Open(old)
	ParseError(err, "While opening "+old)
	defer from.Close()

	to, err2 := os.OpenFile(new, os.O_RDWR|os.O_CREATE, 0666)
	ParseError(err2, "While opening "+new)
	defer to.Close()

	_, err = io.Copy(to, from)
	ParseError(err, "While copying file.")
}

func remove(file string) {
	exec.Command("rm", file).Run()
}

func clean() {

	if target.debug != true && target.clean == true {
		remove("core/ASLR/Mem.map")
		progress()
		remove("core/ASLR/iat/Mem.map")
		progress()
		remove("core/Fixed/Mem.map")
		progress()
		remove("core/Fixed/iat/Mem.map")
		progress()
		remove("stub.o")
		progress()
		remove("Payload")
		progress()
		remove("Payload.key")
		progress()
		exec.Command("sh", "-c", "echo   > stub/payload.h").Run()
		progress()
		exec.Command("sh", "-c", "echo   > stub/key.h").Run()
		progress()
	}
}

func Banner() {

	DIST, _ := exec.Command("lsb_release", "-a").Output()
	if strings.Contains(string(DIST), "Arch") {
		PACKET_MANAGER = "pacman"
		BoldRed.Print(ArchBanner)
	} else {
		BoldRed.Print(BANNER)
	}
	//BoldRed.Print(BANNER)

	BoldBlue.Print("\n# Version: ")
	BoldGreen.Println(VERSION)
	BoldBlue.Print("# Author: ")
	BoldGreen.Println("Ege Balcı")
	BoldBlue.Print("# Source: ")
	BoldGreen.Println("github.com/egebalci/Amber")
}

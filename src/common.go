package main

import "os"
import "fmt"
import "os/exec"
import "strings"
import "gopkg.in/cheggaaa/pb.v1"


func progress() {
	if peid.verbose == false {
		progressBar.Increment()
	}
}

func CreateProgressBar() {

	var full int = 48

	if peid.verbose == false {
		if peid.staged == true {
			full -= 10
		}

		progressBar = pb.New(full)
		progressBar.SetWidth(80)
		progressBar.Start()
	}
}

func verbose(str string, status string) {

	if peid.verbose == true {
		if status == "*" {
			BoldYellow.Print("[*] ")
			white.Println(str)	
		}else if status == "+" {
			BoldGreen.Print("[+] ")
			white.Println(str)	
		}else if status == "-" {
			BoldRed.Print("[-] ")
			white.Println(str)
		}else if status == "!"{
			BoldRed.Print("[!] ")
			white.Println(str)
		}else if status == "Y" {
			BoldYellow.Println(str)
		}else if status == "B" {
			BoldBlue.Println(str)
		}
	}
}

func _verbose(str string, value int32) {

	if peid.verbose == true {
		BoldYellow.Print("[*] ")
		white.Printf(str+" 0x%X\n",value)
	}
}


func CheckRequirements() {

	verbose("Checking requirments...","*")

	if PACKET_MANAGER == "pacman" {
		progress()
		CheckMingw, mingwErr := exec.Command("sh", "-c", "i686-w64-mingw32-g++ --version").Output()
		if !strings.Contains(string(CheckMingw), "Copyright") {
			ParseError(mingwErr,"MingW is not installed.",string(CheckMingw))
		}
		progress()
		CheckNasm, _ := exec.Command("sh", "-c", "nasm -h").Output()
		if !strings.Contains(string(CheckNasm), "usage:") {
			ParseError(nil,"nasm is not installed.",string(CheckNasm))
		}
		progress()
		CheckStrip, _ := exec.Command("sh", "-c", "strip -V").Output()
		if !strings.Contains(string(CheckStrip), "Copyright") {
			ParseError(nil,"strip is not installed.",string(CheckStrip))
		}
		progress()
		CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
		if !strings.Contains(string(CheckXXD), "Amber") {
			ParseError(nil,"xxd is not installed.",string(CheckMingw))
		}
		progress()
		CheckMultiLib, _ := exec.Command("sh", "-c", "pacman -Qs gcc-multilib").Output()
		if strings.Contains(string(CheckMultiLib), "The GNU Compiler") {
			ParseError(nil,"gcc-multilib is not installed.",string(CheckMultiLib))
		}
		progress()
	}else{
		CheckMingw, mingwErr := exec.Command("sh", "-c", "i686-w64-mingw32-g++-win32 --version").Output()
		if !strings.Contains(string(CheckMingw), "Copyright") {
			ParseError(mingwErr,"MingW is not installed.",string(CheckMingw))
		}
		progress()
		CheckNasm, _ := exec.Command("sh", "-c", "nasm -h").Output()
		if !strings.Contains(string(CheckNasm), "usage:") {
			ParseError(nil,"nasm is not installed.",string(CheckNasm))
		}
		progress()
		CheckStrip, _ := exec.Command("sh", "-c", "strip -V").Output()
		if !strings.Contains(string(CheckStrip), "Copyright") {
			ParseError(nil,"strip is not installed.",string(CheckStrip))
		}
		progress()
		CheckXXD, _ := exec.Command("sh", "-c", "echo Amber|xxd").Output()
		if !strings.Contains(string(CheckXXD), "Amber") {
			ParseError(nil,"xxd is not installed.",string(CheckMingw))
		}
		progress()
		CheckMultiLib, _ := exec.Command("sh", "-c", "apt-cache policy gcc-multilib").Output()
		if strings.Contains(string(CheckMultiLib), "(none)") {
			ParseError(nil,"gcc-multilib is not installed.",string(CheckMultiLib))
		}
		progress()
		CheckMultiLibPlus, _ := exec.Command("sh", "-c", "apt-cache policy g++-multilib").Output()
		if strings.Contains(string(CheckMultiLibPlus), "(none)") {
			ParseError(nil,"g++-multilib is not installed.",string(CheckMultiLibPlus))
		}
		progress()
		
	}
}


func ParseError(Err error,ErrStatus string,Msg string){

	if Err != nil {
		
		//progressBar.Finish()
		fmt.Println("\n")
		BoldRed.Println("\n[!] ERROR: "+ErrStatus)
		fmt.Println("\nERROR{\n",Err)
		if len(Msg) > 0 {
			fmt.Println(Msg+"\n}\n")
		}else{
			fmt.Println("\n}\n")
		} 
		if Msg != " " {
			clean()
		}
		fmt.Println("\n")
		os.Exit(1)
	}
}

func Cdir(dir string) {
	err := os.Chdir(dir)
	ParseError(err,"While changing directory.","")
}


func move(Old, New string) {
	err := os.Rename(Old,New)
	if err != nil {
		ParseError(err,"While moving a file.","")
	}
}

func remove(file string) {
	exec.Command("rm", file).Run()
}

func clean() {

	if peid.debug != true {
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
	}else{
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

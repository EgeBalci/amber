package main

import "os"
import "fmt"
import "os/exec"
import "strings"
//import "github.com/fatih/color"
import "gopkg.in/cheggaaa/pb.v1"


func progress() {
	if peid.verbose == false {
		progressBar.Increment()
	}
}

func CreateProgressBar() {

	var full int = 47

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
		fmt.Println(Err)
		if len(Msg) > 0 {
			BoldRed.Println(Msg)
		}
		clean()
		fmt.Println("\n")
		os.Exit(1)
	}
}


func clean() {

	exec.Command("sh", "-c", "rm core/ASLR/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm core/ASLR/iat/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm core/Fixed/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm core/Fixed/iat/Mem.map").Run()
	progress()
	exec.Command("sh", "-c", "rm stub.o").Run()
	progress()
	exec.Command("sh", "-c", "rm Payload").Run()
	progress()
	exec.Command("sh", "-c", "rm Payload.key").Run()
	progress()
	exec.Command("sh", "-c", "echo   > stub/payload.h").Run()
	progress()
	exec.Command("sh", "-c", "echo   > stub/key.h").Run()
	progress()

}

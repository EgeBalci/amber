package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/cheggaaa/pb.v1"
)

func progress() {
	if target.verbose == false {
		progressBar.Increment()
	}
}

func createProgressBar() {
	step := 29
	if target.resource {
		step--
	}
	if target.reflective {
		step = step - 5
	}
	if target.verbose == false {
		progressBar = pb.New(step)
		progressBar.SetWidth(80)
		progressBar.Start()
	}
}

func workdir(fileName string) string {
	defer progress()
	return os.TempDir() + "/" + md5sum(fileName)
}

func md5sum(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	defer progress()
	return hex.EncodeToString(h.Sum(nil))
}

func verbose(str string, status string) {
	if target.verbose {
		switch status {
		case "*":
			BoldYellow.Print("[*] ")
			white.Println(str)
		case "+":
			BoldGreen.Print("[+] ")
			white.Println(str)
		case "-":
			BoldRed.Print("[-] ")
			white.Println(str)
		case "!":
			BoldRed.Print("[!] ")
			white.Println(str)
		case "R":
			BoldRed.Println(str)
		case "Y":
			BoldYellow.Println(str)
		case "B":
			BoldBlue.Println(str)
		case "":
			fmt.Println(str)
		}
	}
}
func print(str string, status string) {

	switch status {
	case "*":
		BoldYellow.Print("[*] ")
		white.Println(str)
	case "+":
		BoldGreen.Print("[+] ")
		white.Println(str)
	case "-":
		BoldRed.Print("[-] ")
		white.Println(str)
	case "!":
		BoldRed.Print("[!] ")
		white.Println(str)
	case "R":
		BoldRed.Println(str)
	case "Y":
		BoldYellow.Println(str)
	case "B":
		BoldBlue.Println(str)
	case "":
		fmt.Println(str)
	}
}

func printParams() {
	BoldYellow.Print("\n[*] File: ")
	BoldBlue.Println(target.FileName)
	BoldYellow.Print("[*] Reflective: ")
	BoldBlue.Println(target.reflective)
	BoldYellow.Print("[*] Key Size: ")
	BoldBlue.Println(target.KeySize)
	BoldYellow.Print("[*] API: ")
	if target.iat {
		BoldBlue.Println("IAT")
	} else {
		BoldBlue.Println("EAT")
	}
	BoldYellow.Print("[*] Verbose: ")
	BoldBlue.Println(target.verbose, "\n")
}

func requirements() {

	verbose("Checking requirments...", "*")

	if runtime.GOOS == "windows" {
		if _, err := os.Stat(os.Getenv("APPDATA") + "\\..\\Local\\bin\\NASM\\nasm.exe"); os.IsNotExist(err) {
			parseErr(errors.New("NASM location could not be detected, nasm is not installed correctly"))
		} else {
			target.nasm = os.Getenv("APPDATA") + "\\..\\Local\\bin\\NASM\\nasm.exe"
		}
	} else {
		if _, err := os.Stat("/usr/bin/nasm"); os.IsNotExist(err) {
			parseErr(errors.New("NASM location could not be detected, nasm is not installed correctly"))
		}
		target.nasm = string("/usr/bin/nasm")
	}
	defer progress()
}

func parseErr(err error) {
	if err != nil {
		if !target.verbose {
			progressBar.Finish()
		}
		fmt.Println("\n")
		clean()
		log.Fatal(err)
	}
}

func fileSize(fileName string) string {
	file, err := os.Open(fileName)
	parseErr(err)
	stat, err := file.Stat()
	parseErr(err)
	defer progress()
	return strconv.Itoa(int(stat.Size()))
}

func cdir(dir string) {
	err := os.Chdir(dir)
	parseErr(err)
	defer progress()
}

func mkdir(dir string) {
	err := os.Mkdir(dir, 0755)
	parseErr(err)
	defer progress()
}

func move(old, new string) {
	err := os.Rename(old, new)
	parseErr(err)
	defer progress()
}

func copyFile(src, dst string) {
	from, err := os.Open(src)
	parseErr(err)
	defer from.Close()
	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	parseErr(err)
	defer to.Close()
	_, err = io.Copy(to, from)
	parseErr(err)
	defer progress()
}

func remove(file string) error {
	defer progress()
	return os.RemoveAll(file)
}

func clean() {
	if !target.debug {
		os.Chdir(os.TempDir())
		os.RemoveAll(target.workdir)
	}
}

func banner() {

	if runtime.GOOS == "windows" {
		print(BasicBanner, "R")
	} else {
		powerlineFonts, err := exec.Command("fc-list").Output()
		parseErr(err)
		if strings.Contains(string(powerlineFonts), "Powerline") {
			print(BANNER, "R")
		} else {
			print(BasicBanner, "R")
		}
	}
	//BoldRed.Print(BANNER)

	BoldBlue.Print("# Version: ")
	BoldGreen.Println(VERSION)
	BoldBlue.Print("# Author: ")
	BoldGreen.Println("Ege Balcı")
	BoldBlue.Print("# Source: ")
	BoldGreen.Println("github.com/egebalci/amber")
}

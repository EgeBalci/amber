package main

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/egebalci/mappe/mape"
	"github.com/rakyll/statik/fs"
)

func compile() {

	mkdir(target.workdir + "/fs")
	mkdir(target.workdir + "/data")
	move(target.workdir+"/stage", target.workdir+"/data/stage")
	id := statik(target.workdir+"/data", target.workdir)

	// Create static fs
	statikFS, err := fs.New()
	parseErr(err)

	// Extract virtual fs libraries
	extract(statikFS, "/stub/fs/fs.go", "./fs/fs.go")
	extract(statikFS, "/stub/fs/walk.go", "./fs/walk.go")
	extract(statikFS, "/stub/fs/fs_test.go", "./fs/fs_test.go")

	// Write stub.go to workdir
	if target.AntiAnalysis {
		extract(statikFS, "/stub/bypass_stub.go", "stub.go")
	} else {
		extract(statikFS, "/stub/stub.go", "stub.go")
	}

	if !target.resource {
		// Write rsrc.syso to workdir
		extract(statikFS, "/stub/rsrc.syso", "rsrc.syso")
	}

	ldflags := `-ldflags=-s`
	if target.subsystem == 0x02 {
		ldflags = `-ldflags=-s -H windowsgui`
	}
	build := exec.Command("go", "build", "-buildmode=exe", ldflags, "-o", "packed.exe", ".")
	build.Env = os.Environ()
	build.Env = append(build.Env, "GOOS=windows")
	if target.arch == "x86" {
		build.Env = append(build.Env, "GOARCH=386")
	} else {
		build.Env = append(build.Env, "GOARCH=amd64")
	}
	obfuscateFunctionNames(id)
	verbose("Compiling go stub...", "*")
	if target.debug {
		build.Stderr = os.Stderr
		build.Stdout = os.Stdout
	}
	err = build.Run()
	parseErr(err)
	copyFile(target.workdir+"/packed.exe", target.FileName)
	if !target.aslr {
		//alignBase(target.FileName)
	}
	defer progress()
}

func obfuscateFunctionNames(id string) {
	verbose("Obfuscating function names...", "*")
	file, err := os.OpenFile(target.workdir+"/stub.go", os.O_WRONLY, os.ModeDevice)
	parseErr(err)
	rawFile, err := ioutil.ReadFile(target.workdir + "/stub.go")
	parseErr(err)
	stub := strings.Replace(string(rawFile), "{{statik}}", id, -1)
	_, err = file.Write([]byte(stub))
	parseErr(err)
	defer progress()
}

func alignBase(fileName string) {

	file, err := pe.Open(fileName)
	parseErr(err)
	opt := mape.ConvertOptionalHeader(file)
	if target.ImageBase != opt.ImageBase {
		verbose("Aligning image base...", "*")
		rawFile, err := ioutil.ReadFile(fileName)
		parseErr(err)
		oFileHeader := int(rawFile[0x3c])

		if target.arch == "x64" {
			oImageBase := int(oFileHeader + 0x30)
			imageBase := binary.LittleEndian.Uint64(rawFile[oImageBase : oImageBase+8])
			verbose("ImageBase "+fmt.Sprintf("0x%016X -> ", imageBase)+fmt.Sprintf("0x%016X", target.ImageBase), "*")
			binary.LittleEndian.PutUint64(rawFile[oImageBase:oImageBase+8], target.ImageBase)
		} else {
			oImageBase := int(oFileHeader + 0x34)
			imageBase := binary.LittleEndian.Uint32(rawFile[oImageBase : oImageBase+4])
			verbose("ImageBase "+fmt.Sprintf("0x%08X -> ", imageBase)+fmt.Sprintf("0x%08X", target.ImageBase), "*")
			binary.LittleEndian.PutUint32(rawFile[oImageBase:oImageBase+4], uint32(target.ImageBase))
		}
		err = ioutil.WriteFile(target.FileName, rawFile, os.ModeDir)
		parseErr(err)

	}

}

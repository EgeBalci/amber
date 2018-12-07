package main

import (
	"os"
	"os/exec"
)

func compile() {
	mkdir(target.workdir + "/data")
	move(target.workdir+"/stage", target.workdir+"/data/stage")
	statik("statik", target.workdir+"/data", target.workdir)

	ldflags := `-ldflags=-s`
	if target.subsystem == 0x02 {
		ldflags = `-ldflags="-s -H windowsgui"`
	}
	build := exec.Command("go", "build", "-buildmode=exe", ldflags, "-o", target.fileName, ".")
	build.Env = os.Environ()
	build.Env = append(build.Env, "GOOS=windows")
	if target.arch == "x86" {
		build.Env = append(build.Env, "GOARCH=386")
	} else {
		build.Env = append(build.Env, "GOARCH=amd64")
	}

	verbose("Compiling go stub...", "*")
	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	err := build.Run()
	parseErr(err)
	defer progress()
}

// func alignBase(fileName string) {
// 	file, err := pe.Open(fileName)
// 	parseErr(err)
// 	opt := mape.ConvertOptionalHeader(file)
// 	if target.imageBase != opt.ImageBase {
// 		rawFile, err := ioutil.ReadFile(fileName)
// 		parseErr(err)
// 		nt := int(rawFile[0x3c])
// 		imageBaseP := int(nt + 0x1c)
// 	}
// }

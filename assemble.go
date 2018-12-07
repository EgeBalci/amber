package main

import (
	"os"
	"os/exec"

	"github.com/egebalci/mappe/mape"
	"github.com/rakyll/statik/fs"
)

func assemble() {

	// Create static fs
	statikFS, err := fs.New()
	parseErr(err)

	Map, err := mape.CreateFileMapping(target.fileName)
	parseErr(err)
	mapFile, err := os.Create(target.workdir + "/Mem.map")
	parseErr(err)
	if target.scrape {
		Map = mape.Scrape(Map)
	} else {
		mapFile.Write(Map)
	}
	mapFile.Close()

	stubName := ""
	apiName := ""
	decoderName := ""

	if target.arch == "x86" {
		if target.aslr == false {
			stubName = "/x86/fixed_stub.asm"
		} else {
			stubName = "/x86/stub.asm"
		}
		if target.iat {
			apiName = "/x86/api/iat_api.asm"
		} else {
			apiName = "/x86/api/block_api.asm"
		}
		decoderName = "/x86/RC4.asm"
	} else {
		if target.aslr == false {
			stubName = "/x64/fixed_stub.asm"

		} else {
			stubName = "/x64/stub.asm"
		}
		if target.iat {
			apiName = "/x64/api/iat_api.asm"
		} else {
			apiName = "/x64/api/block_api.asm"
		}
		decoderName = "/x64/RC4.asm"
	}

	// Write stub.asm to workdir
	stub, err := fs.ReadFile(statikFS, stubName)
	parseErr(err)
	stubFile, err := os.Create(target.workdir + "/stub.asm")
	parseErr(err)
	_, err = stubFile.Write(stub)
	parseErr(err)

	// Write api.asm to workdir
	api, err := fs.ReadFile(statikFS, apiName)
	parseErr(err)
	apiFile, err := os.Create(target.workdir + "/api.asm")
	parseErr(err)
	_, err = apiFile.Write(api)
	parseErr(err)

	// Write RC4.asm to workdir
	decoder, err := fs.ReadFile(statikFS, decoderName)
	parseErr(err)
	decoderFile, err := os.Create(target.workdir + "/RC4.asm")
	parseErr(err)
	_, err = decoderFile.Write(decoder)
	parseErr(err)

	// Write stub.go to workdir
	goStub, err := fs.ReadFile(statikFS, "/stub/stub.go")
	parseErr(err)
	goStubFile, err := os.Create(target.workdir + "/stub.go")
	parseErr(err)
	_, err = goStubFile.Write(goStub)
	parseErr(err)

	if !target.resource {
		// Write rsrc.syso to workdir
		goResource, err := fs.ReadFile(statikFS, "/stub/rsrc.syso")
		parseErr(err)
		goResourceFile, err := os.Create(target.workdir + "/rsrc.syso")
		parseErr(err)
		_, err = goResourceFile.Write(goResource)
		parseErr(err)
	}

	nasm(target.workdir+"/stub.asm", target.workdir+"/payload")
	crypt()
	nasm(target.workdir+"/RC4.asm", target.workdir+"/stage")
	verbose("Assebly completed.", "*")
	defer progress()
}

func nasm(file, out string) {
	err := exec.Command(target.nasm, "-f", "bin", file, "-o", out).Run()
	parseErr(err)
	defer progress()
}

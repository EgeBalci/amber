package main

import (
	"os"
	"os/exec"

	"github.com/egebalci/mappe/mape"
	"github.com/rakyll/statik/fs"
)

func assemble() {
	verbose("Assembling stub...", "*")
	// Create static fs
	statikFS, err := fs.New()
	parseErr(err)

	Map, err := mape.CreateFileMapping(target.FileName)
	parseErr(err)

	err = mape.PerformIntegrityChecks(target.FileName, Map)
	if target.IgnoreIntegrity && err != nil {
		parseErr(err)
	}

	mapFile, err := os.Create(target.workdir + "/Mem.map")
	parseErr(err)
	if target.scrape {
		Map = mape.Scrape(Map)
	} else {
		mapFile.Write(Map)
	}
	mapFile.Close()

	stub := ""
	api := ""
	decoder := ""

	if target.arch == "x86" {
		if target.aslr == false {
			stub = "/x86/fixed_stub.asm"
		} else {
			stub = "/x86/stub.asm"
		}
		if target.iat {
			api = "/x86/api/iat_api.asm"
		} else {
			api = "/x86/api/block_api.asm"
		}
		decoder = "/x86/RC4.asm"
	} else {
		if target.aslr == false {
			stub = "/x64/fixed_stub.asm"

		} else {
			stub = "/x64/stub.asm"
		}
		if target.iat {
			api = "/x64/api/iat_api.asm"
		} else {
			api = "/x64/api/block_api.asm"
		}
		decoder = "/x64/RC4.asm"
	}

	// Write stub.asm to workdir
	extract(statikFS, stub, "stub.asm")

	// Write api.asm to workdir
	extract(statikFS, api, "api.asm")

	// Write RC4.asm to workdir
	extract(statikFS, decoder, "RC4.asm")

	nasm(target.workdir+"/stub.asm", target.workdir+"/payload")
	crypt()
	nasm(target.workdir+"/RC4.asm", target.workdir+"/stage")
	verbose("Assebly completed.", "*")
	defer progress()
}

func nasm(file, out string) {
	nasm := exec.Command(target.nasm, "-f", "bin", file, "-o", out)
	if target.debug {
		nasm.Stderr = os.Stderr
		nasm.Stdout = os.Stdout
	}
	err := nasm.Run()
	parseErr(err)
	defer progress()
}

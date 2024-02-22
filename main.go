package main

import (
	"os"
	"time"

	"github.com/EgeBalci/amber/config"
	amber "github.com/EgeBalci/amber/pkg"
	"github.com/EgeBalci/amber/utils"
	sgn "github.com/EgeBalci/sgn/pkg"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Set globals...
var spinr = spinner.New(spinner.CharSets[9], 30*time.Millisecond)

func main() {
	printBanner()
	cfg, err := config.Parse()
	if err != nil {
		utils.PrintFatal(err)
	}

	encoder, err := sgn.NewEncoder(64)
	if err != nil {
		utils.PrintFatal(err)
	}
	encoder.EncodingCount = cfg.EncodeCount
	encoder.ObfuscationLimit = cfg.ObfuscationLimit
	cfg.PrintSummary()
	// ------------------------------
	pe, err := amber.Open(cfg.FileName)
	if err != nil {
		utils.PrintFatal(err)
	}
	pe.SyscallLoader = cfg.UseSyscalls

	if !pe.HasRelocData {
		utils.PrintErr("%s has no relocation data. Exiting...\n", pe.Name)
		return
		// if pe.ImageBase != 0x400000 {
		// 	utils.PrintErr("Can't switch to fixed address loader because ImageBase mismatch!\n")
		// }
		// utils.PrintStatus("Switching to fixed address loader...\n")
	}

	payload, err := pe.AssembleLoader()
	if err != nil {
		utils.PrintFatal(err)
	}

	if encoder.EncodingCount > 0 {
		spinr.Start()
		spinr.Suffix = " Encoding reflective payload..."
		encoder.SetArchitecture(pe.Architecture)
		payload, err = encoder.Encode(payload)
		if err != nil {
			utils.PrintFatal(err)
		}
		spinr.Stop()
	}

	outFile, err := os.Create(cfg.OutputFile)
	if err != nil {
		utils.PrintFatal(err)
	}

	outFile.Write(payload)
	defer outFile.Close()

	finSize, err := utils.GetFileSize(cfg.OutputFile)
	if err != nil {
		utils.PrintFatal(err)
	}
	utils.PrintStatus("Final Size: %d bytes\n", finSize)
	utils.PrintStatus("Output File: %s\n", cfg.OutputFile)
	utils.PrintGreen("[✔] Reflective PE generated !\n")
}

// BANNER .
const BANNER string = `

//       █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
//      ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
//      ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
//      ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
//      ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
//      ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
//  Reflective PE Packer ☣ Copyright (c) 2017 EGE BALCI
//      %s - %s

`

func printBanner() {
	green := color.New(color.FgGreen).Add(color.Bold)
	red := color.New(color.FgRed).Add(color.Bold)
	blue := color.New(color.FgBlue).Add(color.Bold)
	red.Printf(BANNER, green.Sprintf("v%s", config.Version), blue.Sprintf("https://github.com/egebalci/amber"))
}

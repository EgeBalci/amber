package config

import (
	"fmt"
	"os"

	"github.com/EgeBalci/amber/utils"
	"github.com/alecthomas/kong"
)

const Version = "3.2.0"

func HelpPrompt(options kong.HelpOptions, ctx *kong.Context) error {
	err := kong.DefaultHelpPrinter(options, ctx)
	if err != nil {
		return err
	}
	return nil
}

// Main config struct for parsing the TOML file
type Config struct {
	FileName         string `help:"Input PE file name." name:"file" short:"f"`
	OutputFile       string `help:"Output binary payload file name." name:"out" short:"o"`
	EncodeCount      int    `help:"Number of times to encode the generated reflective payload." name:"encode" short:"e" default:"1"`
	ObfuscationLimit int    `help:"Maximum number of bytes for encoder obfuscation." name:"obfuscate-limit" short:"l" default:"5"`
	UseIAT           bool   `help:"Use IAT API resolver block instead of CRC API resolver block." name:"iat"`
	UseSyscalls      bool   `help:"Perform raw syscalls. (only x64)" name:"sys"`
	ScrapePeHeaders  bool   `help:"Scrape magic byte and DOS stub from PE." name:"scrape"`
	// IgnoreIntegrity  bool   `help:"Ignore PE file integrity check errors." name:"ignore"`
	Verbose bool `help:"Verbose mode." name:"verbose" short:"v"`
	Version kong.VersionFlag
}

// ConfigureOptions accepts a flag set and augments it with agentgo-server
// specific flags. On success, an options structure is returned configured
// based on the selected flags.
func Parse() (*Config, error) {

	cfg := new(Config)
	parser, err := kong.New(
		cfg,
		kong.Help(HelpPrompt),
		kong.UsageOnError(),
		kong.Vars{"version": Version},
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
		}),
	)
	if err != nil {
		return nil, err
	}
	_, err = parser.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if cfg.FileName == "" {
		utils.PrintErr("no file specified! (-f <empty>)\n")
		kong.Help(HelpPrompt)
		os.Exit(1)
	}

	if cfg.OutputFile == "" {
		cfg.OutputFile = fmt.Sprintf("%s.bin", cfg.FileName)
	}

	return cfg, nil
}

func (cfg *Config) PrintSummary() {
	utils.PrintStatus("File: %s\n", cfg.FileName)
	utils.PrintStatus("Encode Count: %d\n", cfg.EncodeCount)
	utils.PrintStatus("Obfuscation Limit: %d\n", cfg.ObfuscationLimit)
	if cfg.UseIAT {
		utils.PrintStatus("API Resolver: IAT\n")
	} else {
		utils.PrintStatus("API Resolver: CRC\n")
	}
	if cfg.UseSyscalls {
		utils.PrintStatus("Raw Syscalls: True\n")
	}
}

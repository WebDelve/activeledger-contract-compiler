package main

import (
	"flag"

	"github.com/WebDelve/activeledger-contract-compiler/compiler"
	"github.com/WebDelve/activeledger-contract-compiler/config"
)

type CLIFlags struct {
	contractEntry  string
	contractOutput string
}

func main() {
	config := config.GetConfig()
	flags := getFlags()

	var blank string
	if flags.contractOutput != blank {
		config.Output = flags.contractOutput
	}

	comp := compiler.GetCompiler(config, flags.contractEntry)
	comp.Compile()
}

func getFlags() CLIFlags {
	contractEntryPtr := flag.String(
		"p",
		"smartcontract/main.ts",
		"Path to Smart Contract entry point",
	)

	contractOutPtr := flag.String(
		"o",
		"",
		"Output file",
	)

	flag.Parse()

	flags := CLIFlags{
		contractEntry:  *contractEntryPtr,
		contractOutput: *contractOutPtr,
	}

	return flags
}

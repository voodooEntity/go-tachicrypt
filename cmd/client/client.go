package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/voodooEntity/go-tachicrypt/src/core"
	"github.com/voodooEntity/go-tachicrypt/src/prettywriter"
	"github.com/voodooEntity/go-tachicrypt/src/utils"
)

var version = "Alpha 0.2.0"

func main() {
	// Define flags
	hide := flag.Bool("hide", false, "Hide (encrypt) data")
	unhide := flag.Bool("unhide", false, "Unhide (decrypt) data")
	dataPath := flag.String("data", "", "Path to the data file or directory")
	partCount := flag.Int("parts", -1, "Amount of parts that should be created")
	outputDir := flag.String("output", "", "Output directory for encrypted data or decrypted data")
	help := flag.Bool("help", false, "Show help message")

	// Parse flags
	flag.Parse()

	// Show help if --help is specified
	if *help {
		printUsage()
		return
	}

	if *hide && -1 == *partCount {
		utils.ExitError("Missing mandatory --parts parameter \n")
		return
	}

	if *hide && *unhide {
		utils.ExitError("Cannot use both --hide and --unhide options at the same time. \n")
		return
	}

	if (*hide || *unhide) && (*dataPath == "" || *outputDir == "") {
		utils.ExitError("Both --data and --output must be specified. \n")
		return
	}

	utils.PrintApplicationHeader(version)

	if *hide {
		c := core.New()
		prefilledPwd := os.Getenv("TACHICRYPT_PASSWORD")
		err := c.Hide(*dataPath, *partCount, *outputDir, prefilledPwd)
		if err != nil {
			utils.ExitError(fmt.Sprintf("Error hiding data: %v \n", err))
			return
		}
		return
	}

	if *unhide {
		c := core.New()
		prefilledPwd := os.Getenv("TACHICRYPT_PASSWORD")
		err := c.Unhide(*dataPath, *outputDir, prefilledPwd)
		if err != nil {
			utils.ExitError(fmt.Sprintf("Error unhiding data: %v \n", err))
			return
		}
		return
	}
	printUsage()
}

// printUsage prints the usage information for the command-line tool.
func printUsage() {
	prettywriter.WriteInBox(40, "Usage: tachicrypt [options]", prettywriter.Green, prettywriter.BlackBG, prettywriter.DoubleLine)
	fmt.Println("")
	prettywriter.Writeln("Options:", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --hide             Hide (encrypt) data", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --unhide           Unhide (decrypt) data", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --data     [arg]   Path to the data file or directory", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --parts    [arg]   Amount of parts to be created when hiding", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --output   [arg]   Output directory for encrypted data or decrypted data", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  --help             Show this help message", prettywriter.Green, prettywriter.BlackBG)
	fmt.Println("")
	prettywriter.Writeln("Examples:", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  Encrypt data: tachicrypt --hide --parts 10 --data /path/to/data --output /path/to/output", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("  Decrypt data: tachicrypt --data /path/to/encrypted/data --unhide --output /path/to/output ", prettywriter.Green, prettywriter.BlackBG)
	fmt.Println("")
}

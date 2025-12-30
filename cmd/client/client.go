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

// test hook for ExitError to avoid os.Exit in unit tests; defaults to utils.ExitError
var exitErrorFn = utils.ExitError

// test hooks to allow stubbing core operations in unit tests
var hideFunc = func(dataPath string, partCount int, outputDir string, prefilledPassword string) error {
    return core.New().Hide(dataPath, partCount, outputDir, prefilledPassword)
}
var unhideFunc = func(dataPath string, outputDir string, prefilledPassword string) error {
    return core.New().Unhide(dataPath, outputDir, prefilledPassword)
}

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

 // Validate flags; on failure exitErrorFn will be invoked and we return
 if !validateFlags(*hide, *unhide, *partCount, *dataPath, *outputDir) {
     return
 }

	utils.PrintApplicationHeader(version)

 if *hide {
        prefilledPwd := os.Getenv("TACHICRYPT_PASSWORD")
        err := hideFunc(*dataPath, *partCount, *outputDir, prefilledPwd)
        if err != nil {
            exitErrorFn(fmt.Sprintf("Error hiding data: %v \n", err))
        }
        return
    }

 if *unhide {
        prefilledPwd := os.Getenv("TACHICRYPT_PASSWORD")
        err := unhideFunc(*dataPath, *outputDir, prefilledPwd)
        if err != nil {
            exitErrorFn(fmt.Sprintf("Error unhiding data: %v \n", err))
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

// validateFlags performs CLI flags validation and invokes exitErrorFn on failure.
// Returns true if validation succeeded and execution can continue.
func validateFlags(hide, unhide bool, parts int, dataPath, outputDir string) bool {
    if hide && -1 == parts {
        exitErrorFn("Missing mandatory --parts parameter \n")
        return false
    }
    if hide && unhide {
        exitErrorFn("Cannot use both --hide and --unhide options at the same time. \n")
        return false
    }
    if (hide || unhide) && (dataPath == "" || outputDir == "") {
        exitErrorFn("Both --data and --output must be specified. \n")
        return false
    }
    return true
}

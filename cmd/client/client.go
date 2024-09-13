package main

import (
	"flag"
	"fmt"
	"github.com/voodooEntity/go-tachicrypt/src/core"
	"log"
)

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
		log.Fatal("Missing mandatory --parts parameter")
	}

	if *hide && *unhide {
		log.Fatal("Cannot use both --hide and --unhide options at the same time.")
	}

	if (*hide || *unhide) && (*dataPath == "" || *outputDir == "") {
		log.Fatal("Both --data and --output must be specified.")
	}

	if *hide {
		c := core.New()
		err := c.Hide(*dataPath, *partCount, *outputDir, "")
		if err != nil {
			log.Fatalf("Error hiding data: %v", err)
		}
		log.Println("Data encrypted and saved successfully.")
	}

	if *unhide {
		c := core.New()
		err := c.Unhide(*dataPath, *outputDir, "")
		if err != nil {
			log.Fatalf("Error unhiding data: %v", err)
		}
		log.Println("Data decrypted and saved successfully.")
	}
}

// printUsage prints the usage information for the command-line tool.
func printUsage() {
	fmt.Println("Usage: tachicrypt [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --hide           Hide (encrypt) data")
	fmt.Println("  --unhide         Unhide (decrypt) data")
	fmt.Println("  --data           Path to the data file or directory")
	fmt.Println("  --parts          Amount of parts to be created when hiding")
	fmt.Println("  --output         Output directory for encrypted data or decrypted data")
	fmt.Println("  --help           Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  Encrypt data: tachicrypt --hide --parts 10 --data /path/to/data --output /path/to/output")
	fmt.Println("  Decrypt data: tachicrypt --unhide --output /path/to/output ")
}

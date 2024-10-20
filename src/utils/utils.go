package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/voodooEntity/go-tachicrypt/src/prettywriter"
	"golang.org/x/term"
	"os"
	"syscall"
)

var TachiHeading = `
 _____          _     _ _____                  _   
|_   _|        | |   (_)  __ \                | |  
  | | __ _  ___| |__  _| /  \/_ __ _   _ _ __ | |_ 
  | |/ _  |/ __| '_ \| | |   | '__| | | | '_ \| __|
  | | (_| | (__| | | | | \__/\ |  | |_| | |_) | |_ 
  \_/\__,_|\___|_| |_|_|\____/_|   \__, | .__/ \__|
                                    __/ | |        
                                   |___/|_|        
`

func GenerateRandomFilename() (string, error) {
	filename := make([]byte, 16) // Adjust length as needed
	if _, err := rand.Read(filename); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", filename), nil
}

func PromptForPassword(prompt string) string {
	prettywriter.WriteInBox(40, prompt, prettywriter.BlackBG, prettywriter.Green, prettywriter.DoubleLine)
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		prettywriter.Writeln(fmt.Sprintf("\nError reading password: %+v", err), prettywriter.Red, prettywriter.BlackBG)
		fmt.Sprintf("\nError reading password: %+v", err)
		return PromptForPassword(prompt)
	}
	if string(password) == "" {
		prettywriter.Writeln("Error: Invalid empty password.", prettywriter.Red, prettywriter.BlackBG)
		return PromptForPassword(prompt)
	}
	prettywriter.Write("[**] Password entered successfully", prettywriter.BlackBG, prettywriter.Green)
	fmt.Println("")
	pwd := string(password)
	return pwd
}

func PrintApplicationHeader(version string) {
	prettywriter.Write(TachiHeading, prettywriter.BlackBG, prettywriter.Green)
	prettywriter.Writeln("> Version: "+version, prettywriter.BlackBG, prettywriter.Green)
	prettywriter.Writeln("> Github: https://github.com/voodooEntity/go-tachicrypt", prettywriter.BlackBG, prettywriter.Green)
	prettywriter.Writeln("> Author: voodooEntity", prettywriter.BlackBG, prettywriter.Green)
	prettywriter.Writeln("> Contributors: f0o", prettywriter.BlackBG, prettywriter.Green)
	fmt.Println("")
}

func ExitError(message string) {
	fmt.Println("")
	prettywriter.WriteInBox(40, "!Error!", prettywriter.Black, prettywriter.Red, prettywriter.DoubleLine)
	fmt.Println("")
	prettywriter.Writeln(message, prettywriter.Black, prettywriter.Red)
	fmt.Println("")
	os.Exit(1)
}

func ConcatByteSlices(byteSlices [][]byte) []byte {
	totalLength := 0
	for _, slice := range byteSlices {
		totalLength += len(slice)
	}
	result := make([]byte, totalLength)
	offset := 0
	for _, slice := range byteSlices {
		copy(result[offset:], slice)
		offset += len(slice)
	}
	return result
}

package utils

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/term"
	"syscall"
)

func GenerateRandomFilename() (string, error) {
	filename := make([]byte, 16) // Adjust length as needed
	if _, err := rand.Read(filename); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", filename), nil
}

func PromptForPassword(prompt string) string {
	fmt.Println("\n" + prompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		return PromptForPassword(prompt)
	}

	fmt.Println("\nPassword entered successfully.")
	pwd := string(password) // Use the password as needed
	return pwd
}

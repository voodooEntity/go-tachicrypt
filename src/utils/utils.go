package utils

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
)

// GenerateRandomFilename generates a random filename.
func GenerateRandomFilename() (string, error) {
	filename := make([]byte, 16) // Adjust length as needed
	if _, err := rand.Read(filename); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", filename), nil
}

// PromptForPassword prompts the user for a password securely.
func PromptForPassword(message string) (string, error) {
	fmt.Println("\n" + message)
	password, err := readPassword()
	if err != nil {
		return "", err
	}
	return password, nil
}

// readPassword reads a password from standard input without echoing it.
func readPassword() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	pwd, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return "", err
	}
	return pwd, nil
}

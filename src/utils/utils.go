package utils

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"syscall"
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
func PromptForPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := readPassword()
	if err != nil {
		return "", err
	}
	return password, nil
}

// readPassword reads a password from standard input without echoing it.
func readPassword() (string, error) {
	// Disable echo on Unix-like systems
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 0, Max: 0}); err != nil {
		return "", err
	}
	defer syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 0, Max: 0})

	// Read password
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	password := scanner.Text()

	return password, scanner.Err()
}

package fileutils

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func WriteToFile(filename string, content []byte) error {
	// Open the file for writing with 0644 permissions
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close() // Ensure the file is closed even if an error occurs

	// Write the content to the file
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func ObfuscateFileTimestamps(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		if err := touchFile(filePath); err != nil {
			return err
		}
	}
	return nil
}

func touchFile(filePath string) error {
	rand.Seed(time.Now().UnixNano())
	randomTime := time.Now().Add(time.Duration(rand.Intn(1000000)-500000) * time.Second)
	return os.Chtimes(filePath, randomTime, randomTime)
}

package fileutils

import (
	"math/rand"
	"os"
	"time"
)

func WriteToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
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

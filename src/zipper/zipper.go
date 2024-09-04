package zipper

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

// Zipper struct that holds no state, but groups functionality together.
type Zipper struct{}

// New creates a new Zipper instance.
func New() *Zipper {
	return &Zipper{}
}

// ZipToString takes a path (file or directory) and returns a base64-encoded string of the zipped content.
func (z *Zipper) ZipToString(path string) (string, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	err := z.addFiles(zipWriter, path, "")
	if err != nil {
		return "", err
	}

	err = zipWriter.Close()
	if err != nil {
		return "", err
	}

	// Convert the buffer to a base64-encoded string
	zippedBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return zippedBase64, nil
}

// UnzipFromString takes a base64-encoded zip string and unpacks it to the specified destination path.
func (z *Zipper) UnzipFromString(encodedZip, destPath string) error {
	// Decode the base64 string
	zipData, err := base64.StdEncoding.DecodeString(encodedZip)
	if err != nil {
		return err
	}

	// Create a reader for the zip data
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	// Extract files from the zip archive
	for _, file := range zipReader.File {
		filePath := filepath.Join(destPath, file.Name)

		// If it's a directory, create it
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// If it's a file, write it to the disk
		if err := z.extractFile(file, filePath); err != nil {
			return err
		}
	}

	return nil
}

// addFiles recursively adds files to the zip archive.
func (z *Zipper) addFiles(zipWriter *zip.Writer, basePath, baseInZip string) error {
	info, err := os.Stat(basePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		files, err := os.ReadDir(basePath)
		if err != nil {
			return err
		}
		for _, file := range files {
			fPath := filepath.Join(basePath, file.Name())
			if err := z.addFiles(zipWriter, fPath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		}
	} else {
		fileToZip, err := os.Open(basePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Create a new file in the zip archive
		zipFile, err := zipWriter.Create(baseInZip)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, fileToZip)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractFile writes the contents of a zipped file to the specified destination on disk.
func (z *Zipper) extractFile(file *zip.File, destPath string) error {
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Open the file in the archive for reading
	fileInArchive, err := file.Open()
	if err != nil {
		return err
	}
	defer fileInArchive.Close()

	_, err = io.Copy(destFile, fileInArchive)
	return err
}

// UnzipFromBytes extracts a ZIP archive from a byte slice to the specified directory.
func (z *Zipper) UnzipFromBytes(zipData []byte, outputDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		err := extractFile(file, outputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts a single file from the ZIP archive.
func extractFile(file *zip.File, outputDir string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create the file path
	fpath := filepath.Join(outputDir, file.Name)

	// Check if the file is a directory
	if file.FileInfo().IsDir() {
		return os.MkdirAll(fpath, file.Mode())
	}

	// Create the file
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the file content
	_, err = io.Copy(f, rc)
	if err != nil {
		return err
	}

	return nil
}

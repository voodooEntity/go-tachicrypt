package zipper

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

type Zipper struct {
}

func New() *Zipper {
	return &Zipper{}
}

func (z *Zipper) ZipToString(path string) (string, error) {
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)

	// Zip the file
	err := z.zipFile(path, "", w)
	if err != nil {
		w.Close() // Ensure to close the writer in case of error
		return "", err
	}

	// Close the writer before encoding
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode the zip file as base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (z *Zipper) zipFile(path string, prefix string, w *zip.Writer) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		prefix = filepath.Join(prefix, filepath.Base(path))
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			if err := z.zipFile(filepath.Join(path, file.Name()), prefix, w); err != nil {
				return err
			}
		}
	} else {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		zf, err := w.Create(filepath.Join(prefix, filepath.Base(path)))
		if err != nil {
			return err
		}
		_, err = io.Copy(zf, f)
		return err
	}

	return nil
}

func (z *Zipper) ExtractBase64ZipToDir(zipBase64 string, destDir string) error {
	r, err := base64.StdEncoding.DecodeString(zipBase64)
	if err != nil {
		return err
	}

	// Create a ZIP reader directly from the byte slice
	reader, err := zip.NewReader(bytes.NewReader(r), int64(len(r)))
	if nil != err {
		return err
	}

	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filepath.Join(destDir, f.Name), f.FileInfo().Mode())
			if err != nil {
				return err
			}
		} else {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			outFile := filepath.Join(destDir, f.Name)

			err = os.MkdirAll(filepath.Dir(outFile), 0755) // ### note to myself - 755 fixes a current issue but is not a good chmod
			if err != nil {
				return err
			}

			fw, err := os.Create(outFile)
			if err != nil {
				return err
			}

			_, err = io.Copy(fw, rc)
			if err != nil {
				return err
			}
			fw.Close()
			rc.Close()
		}
	}

	return nil
}

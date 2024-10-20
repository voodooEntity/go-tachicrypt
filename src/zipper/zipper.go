package zipper

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type Zipper struct {
}

func New() *Zipper {
	return &Zipper{}
}

func (z *Zipper) Zip(path string) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := zip.NewWriter(buf)

	// Zip the file(s)
	err := z.zipFile(path, "", w)
	if err != nil {
		w.Close()
		return []byte{}, err
	}

	// Close the writer before encoding
	err = w.Close()
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
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

func (z *Zipper) Extract(zipData []byte, destDir string) error {
	// Create a ZIP reader directly from the byte slice
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
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

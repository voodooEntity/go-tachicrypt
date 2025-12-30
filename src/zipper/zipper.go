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

// test hooks for easier unit testing of error paths; default to real implementations
var (
    osStatFn       = os.Stat
    osReadDirFn    = os.ReadDir
    osOpenFn       = os.Open
    osMkdirAllFn   = os.MkdirAll
    osCreateFn     = os.Create
    ioCopyFn       = io.Copy
    zipNewReaderFn = func(b []byte) (*zip.Reader, error) { return zip.NewReader(bytes.NewReader(b), int64(len(b))) }
    zipFileOpenFn  = func(f *zip.File) (io.ReadCloser, error) { return f.Open() }
    closeZipWriterFn = func(w *zip.Writer) error { return w.Close() }
    createZipEntryFn = func(w *zip.Writer, name string) (io.Writer, error) { return w.Create(name) }
)

func (z *Zipper) Zip(path string) ([]byte, error) {
    buf := &bytes.Buffer{}
    w := zip.NewWriter(buf)

    // Zip the file(s)
    err := z.zipFile(path, "", w)
    if err != nil {
        _ = closeZipWriterFn(w)
        return []byte{}, err
    }

    // Close the writer before encoding
    err = closeZipWriterFn(w)
    if err != nil {
        return []byte{}, err
    }

    return buf.Bytes(), nil
}

func (z *Zipper) zipFile(path string, prefix string, w *zip.Writer) error {
    info, err := osStatFn(path)
    if err != nil {
        return err
    }

    if info.IsDir() {
        prefix = filepath.Join(prefix, filepath.Base(path))
        files, err := osReadDirFn(path)
        if err != nil {
            return err
        }
        for _, file := range files {
            if err := z.zipFile(filepath.Join(path, file.Name()), prefix, w); err != nil {
                return err
            }
        }
    } else {
        f, err := osOpenFn(path)
        if err != nil {
            return err
        }
        defer f.Close()

        zf, err := createZipEntryFn(w, filepath.Join(prefix, filepath.Base(path)))
        if err != nil {
            return err
        }
        _, err = ioCopyFn(zf, f)
        return err
    }

    return nil
}

func (z *Zipper) Extract(zipData []byte, destDir string) error {
    // Create a ZIP reader directly from the byte slice
    reader, err := zipNewReaderFn(zipData)
    if nil != err {
        return err
    }

    for _, f := range reader.File {
        if f.FileInfo().IsDir() {
            err := osMkdirAllFn(filepath.Join(destDir, f.Name), f.FileInfo().Mode())
            if err != nil {
                return err
            }
        } else {
            rc, err := zipFileOpenFn(f)
            if err != nil {
                return err
            }

            outFile := filepath.Join(destDir, f.Name)

            err = osMkdirAllFn(filepath.Dir(outFile), 0755) // ### note to myself - 755 fixes a current issue but is not a good chmod
            if err != nil {
                return err
            }

            fw, err := osCreateFn(outFile)
            if err != nil {
                return err
            }

            _, err = ioCopyFn(fw, rc)
            if err != nil {
                return err
            }
            fw.Close()
            rc.Close()
        }
    }

    return nil
}

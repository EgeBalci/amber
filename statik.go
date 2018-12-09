package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rakyll/statik/fs"
)

var namePackage string = "statik"

// mtimeDate holds the arbitrary mtime that we assign to files when
// flagNoMtime is set.
var mtimeDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func statik(src, dest string) string {
	verbose("Creating virtual file system...", "*")
	file, id, err := generateSource(src)
	parseErr(err)
	destDir := path.Join(dest, namePackage)
	err = os.MkdirAll(destDir, 0755)
	parseErr(err)
	err = rename(file.Name(), path.Join(destDir, RandomString(10)+".go"))
	parseErr(err)
	defer progress()
	return id
}

func extract(hfs http.FileSystem, fileName string, dst string) {
	rawFile, err := fs.ReadFile(hfs, fileName)
	parseErr(err)
	file, err := os.Create(target.workdir + "/" + dst)
	parseErr(err)
	defer file.Close()
	_, err = file.Write(rawFile)
	parseErr(err)
	defer progress()
}

// rename tries to os.Rename, but fall backs to copying from src
// to dest and unlink the source if os.Rename fails.
func rename(src, dest string) error {
	// Try to rename generated source.
	if err := os.Rename(src, dest); err == nil {
		return nil
	}
	// If the rename failed (might do so due to temporary file residing on a
	// different device), try to copy byte by byte.
	rc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
		os.Remove(src) // ignore the error, source is in tmp.
	}()

	if _, err = os.Stat(dest); !os.IsNotExist(err) {
		return errors.New("file " + dest + " already exists; use -f to overwrite")
	}

	wc, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer wc.Close()

	if _, err = io.Copy(wc, rc); err != nil {
		// Delete remains of failed copy attempt.
		os.Remove(dest)
	}
	return err
}

// Walks on the source path and generates source code
// that contains source directory's contents as zip contents.
// Generates source registers generated zip contents data to
// be read by the statik/fs HTTP file system.
func generateSource(srcPath string) (file *os.File, id string, err error) {
	var (
		buffer    bytes.Buffer
		zipWriter io.Writer
	)

	zipWriter = &buffer
	f, err := ioutil.TempFile("", namePackage)
	if err != nil {
		return
	}

	zipWriter = io.MultiWriter(zipWriter, f)
	defer f.Close()

	w := zip.NewWriter(zipWriter)
	if err = filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories and hidden files.
		// No entry is needed for directories in a zip file.
		// Each file is represented with a path, no directory
		// entities are required to build the hierarchy.
		if fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
			return nil
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fHeader, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fHeader.Name = filepath.ToSlash(relPath)
		fHeader.Method = zip.Deflate
		f, err := w.CreateHeader(fHeader)
		if err != nil {
			return err
		}
		_, err = f.Write(b)
		return err
	}); err != nil {
		return nil, "", err
	}
	if err = w.Close(); err != nil {
		return nil, "", err
	}
	funcID := RandomString(10)
	// then embed it as a quoted string
	var qb bytes.Buffer
	fmt.Fprintf(&qb, `
package %s

import (
	"../fs"
)

func `+funcID+`() {
	data := "`, namePackage)
	FprintZipData(&qb, buffer.Bytes())
	fmt.Fprint(&qb, `"
	fs.Register(data)
}
`)

	if err = ioutil.WriteFile(f.Name(), qb.Bytes(), 0644); err != nil {
		return nil, "", err
	}
	return f, funcID, nil
}

// FprintZipData converts zip binary contents to a string literal.
func FprintZipData(dest *bytes.Buffer, zipData []byte) {
	for _, b := range zipData {
		if b == '\n' {
			dest.WriteString(`\n`)
			continue
		}
		if b == '\\' {
			dest.WriteString(`\\`)
			continue
		}
		if b == '"' {
			dest.WriteString(`\"`)
			continue
		}
		if (b >= 32 && b <= 126) || b == '\t' {
			dest.WriteByte(b)
			continue
		}
		fmt.Fprintf(dest, "\\x%02x", b)
	}
}

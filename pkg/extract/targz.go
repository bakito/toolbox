package extract

import (
	"archive/tar"
	"archive/zip"
	"github.com/verybluebot/tarinator-go"
	"github.com/xi2/xz"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func File(file, target string) error {
	if strings.HasSuffix(file, ".tar.gz") {
		return tarGz(file, target)
	}
	if strings.HasSuffix(file, ".zip") {
		return unzip(file, target)
	}
	if strings.HasSuffix(file, ".tar.xz") {
		return tarXz(file, target)
	}
	return nil
}

func unzip(file string, target string) error {
	read, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(target, file.Name)
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func tarGz(file, target string) error {
	return tarinator.UnTarinate(target, file)
}

func tarXz(file, target string) error {
	// Open a file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	// Create an xz Reader
	r, err := xz.NewReader(f, 0)
	if err != nil {
		return err
	}
	// Create a tar Reader
	tr := tar.NewReader(r)
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			// create a directory
			err = os.MkdirAll(filepath.Join(target, hdr.Name), 0777)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// write a file
			w, err := os.Create(filepath.Join(target, hdr.Name))
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tr)
			if err != nil {
				return err
			}
			_ = w.Close()
		}
	}
	return nil
}

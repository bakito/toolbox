package extract

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/verybluebot/tarinator-go"
	"github.com/xi2/xz"
)

func File(file, target string) (bool, error) {
	if strings.HasSuffix(file, ".tar.gz") {
		log.Printf("Extracting %s", file)
		return true, tarGz(file, target)
	}
	if strings.HasSuffix(file, ".zip") {
		log.Printf("Extracting %s", file)
		return true, unzip(file, target)
	}
	if strings.HasSuffix(file, ".tar.xz") {
		log.Printf("Extracting %s", file)
		return true, tarXz(file, target)
	}
	return false, nil
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
		// #nosec G305
		name := path.Join(target, file.Name)
		_ = os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		_, _ = create.ReadFrom(open)
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
	// Create a xz Reader
	r, err := xz.NewReader(f, 0)
	if err != nil {
		return err
	}
	// Create a tar Reader
	tr := tar.NewReader(r)
	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			// create a directory
			// #nosec G305
			err = os.MkdirAll(filepath.Join(target, hdr.Name), 0o777)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// write a file
			// #nosec G305
			path := filepath.Join(target, hdr.Name)
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			// #nosec G110
			_, err = io.Copy(w, tr)
			if err != nil {
				return err
			}
			_ = w.Close()
		}
	}
	return nil
}

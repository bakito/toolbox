package extract

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bakito/toolbox/pkg/quietly"
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
	defer quietly.Close(read)
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		if err := unzipFile(file, target); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, target string) error {
	open, err := file.Open()
	if err != nil {
		return err
	}
	name, err := sanitizeArchivePath(target, file.Name)
	if err != nil {
		return err
	}
	parent, _ := filepath.Split(name)
	_ = os.MkdirAll(parent, os.ModeDir)
	create, err := os.Create(name)
	if err != nil {
		return err
	}
	defer quietly.Close(create)
	_, err = create.ReadFrom(open)
	return err
}

// sanitize archive file pathing from "G305: Zip Slip vulnerability"
func sanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", t)
}

func tarGz(file, target string) error {
	tarFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer quietly.Close(tarFile)

	uncompressedStream, err := gzip.NewReader(tarFile)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("extractTarGz: Next() failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			//nolint:gosec
			if err := os.MkdirAll(filepath.Join(target, header.Name), 0o755); err != nil {
				return fmt.Errorf("extractTarGz: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			if err := extractTarFile(target, header, tarReader); err != nil {
				return err
			}

		default:
			return fmt.Errorf("extractTarGz: uknown type: %d in %s",
				header.Typeflag,
				header.Name)
		}
	}
	return nil
}

func extractTarFile(target string, header *tar.Header, tarReader *tar.Reader) error {
	//nolint:gosec
	path := filepath.Join(target, header.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("extractTarGz: Mkdir() failed: %w", err)
	}
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("extractTarGz: Create() failed: %w", err)
	}
	defer quietly.Close(outFile)
	if _, err := io.Copy(outFile, tarReader); err != nil {
		return fmt.Errorf("extractTarGz: Copy() failed: %w", err)
	}
	return nil
}

func tarXz(file, target string) error {
	// Open a file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer quietly.Close(f)
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

		if err := tarXzFile(tr, hdr, target); err != nil {
			return err
		}

	}
	return nil
}

func tarXzFile(tr *tar.Reader, hdr *tar.Header, target string) error {
	path, err := sanitizeArchivePath(target, hdr.Name)
	if err != nil {
		return err
	}
	switch hdr.Typeflag {
	case tar.TypeDir:
		// create a directory
		err = os.MkdirAll(path, 0o777)
		if err != nil {
			return err
		}
	case tar.TypeReg:
		// write a file
		w, err := os.Create(path)
		if err != nil {
			return err
		}

		defer quietly.Close(w)
		for {
			_, err := io.CopyN(w, tr, 1024)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
		}

	}
	return nil
}

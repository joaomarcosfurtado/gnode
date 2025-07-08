package extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Extractor struct{}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) ExtractTarGz(reader io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("error creating reader gzip: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading header from tar: %v", err)
		}

		pathParts := strings.Split(header.Name, "/")
		if len(pathParts) <= 1 {
			continue
		}

		relativePath := strings.Join(pathParts[1:], "/")
		if relativePath == "" {
			continue
		}

		path := filepath.Join(destDir, relativePath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
		case tar.TypeReg:
			if err := e.extractFile(tr, path, header.Mode); err != nil {
				return fmt.Errorf("error extracting file %v", err)
			}
		}
	}

	return nil
}

func (e *Extractor) extractFile(tr *tar.Reader, path string, mode int64) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(mode))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, tr); err != nil {
		return err
	}

	return nil
}

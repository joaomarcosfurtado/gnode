package extractor

import (
	"archive/tar"
	"archive/zip"
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

func (e *Extractor) ExtractZip(src, destDir string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("error opening zip file: %v", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		pathParts := strings.Split(file.Name, "/")
		if len(pathParts) <= 1 {
			continue
		}

		relativePath := strings.Join(pathParts[1:], "/")
		if relativePath == "" {
			continue
		}

		path := filepath.Join(destDir, relativePath)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
			continue
		}

		if err := e.extractZipFile(file, path); err != nil {
			return fmt.Errorf("error extracting file %s: %v", file.Name, err)
		}
	}

	return nil
}

func (e *Extractor) extractZipFile(file *zip.File, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
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

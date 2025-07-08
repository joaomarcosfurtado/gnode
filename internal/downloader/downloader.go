package downloader

import (
	"fmt"
	"io"
	"net/http"
)

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(url string) (io.ReadCloser, error) {
	fmt.Printf("Downloading from %s...\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading: %v", err)
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("error downloading: status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

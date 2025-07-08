package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NodeVersion struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

type Service struct {
	baseURL string
}

func NewService(baseUrl string) *Service {
	return &Service{
		baseURL: baseUrl,
	}
}

func (s *Service) ListRemote() ([]NodeVersion, error) {
	resp, err := http.Get(s.baseURL + "/index.json")

	if err != nil {
		return nil, fmt.Errorf("error getting remote versions: %v", err)
	}
	defer resp.Body.Close()

	var versions []NodeVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("error decoding JSON %v", err)
	}

	return versions, nil
}

func (s *Service) NormalizeVersion(version string) string {
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}

func (s *Service) GetDownloadURL(version, goos, goarch string) string {
	filename := fmt.Sprintf("node-%s-%s-%s.tar.gz", s.baseURL, version, goarch)
	return fmt.Sprintf("%s/%s/%s", s.baseURL, version, filename)
}

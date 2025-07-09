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
	platform := ""
	switch goos {
	case "windows":
		platform = "win"
	case "darwin":
		platform = "darwin"
	default:
		platform = goos
	}

	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	if platform == "win" {
		filename := fmt.Sprintf("node-%s-%s-%s.zip", version, platform, arch)
		return fmt.Sprintf("%s/%s/%s", s.baseURL, version, filename)
	}

	ext := ".tar.gz"
	filename := fmt.Sprintf("node-%s-%s-%s%s", version, platform, arch, ext)
	return fmt.Sprintf("%s/%s/%s", s.baseURL, version, filename)
}

func (s *Service) CheckAvailableFiles(version, goarch string) map[string]bool {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	urls := map[string]string{
		"zip":     fmt.Sprintf("%s/%s/node-%s-win-%s.zip", s.baseURL, version, version, arch),
		"node":    fmt.Sprintf("%s/%s/win-%s/node.exe", s.baseURL, version, arch),
		"npm":     fmt.Sprintf("%s/%s/win-%s/npm", s.baseURL, version, arch),
		"npm.cmd": fmt.Sprintf("%s/%s/win-%s/npm.cmd", s.baseURL, version, arch),
		"npx":     fmt.Sprintf("%s/%s/win-%s/npx", s.baseURL, version, arch),
		"npx.cmd": fmt.Sprintf("%s/%s/win-%s/npx.cmd", s.baseURL, version, arch),
	}

	available := make(map[string]bool)

	for name, url := range urls {
		resp, err := http.Head(url)
		if err == nil && resp.StatusCode == 200 {
			available[name] = true
			resp.Body.Close()
		} else {
			available[name] = false
		}
	}

	return available
}

func (s *Service) GetDownloadStrategy(version, goarch string) string {
	available := s.CheckAvailableFiles(version, goarch)

	if available["zip"] {
		return "zip"
	} else if available["node"] {
		return "binaries"
	} else {
		return "unknown"
	}
}

func (s *Service) GetWindowsNodeURL(version, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	return fmt.Sprintf("%s/%s/win-%s/node.exe", s.baseURL, version, arch)
}

func (s *Service) GetWindowsNpmURL(version, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	return fmt.Sprintf("%s/%s/win-%s/npm", s.baseURL, version, arch)
}

func (s *Service) GetWindowsNpmCmdURL(version, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	return fmt.Sprintf("%s/%s/win-%s/npm.cmd", s.baseURL, version, arch)
}

func (s *Service) GetWindowsNpxURL(version, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	return fmt.Sprintf("%s/%s/win-%s/npx", s.baseURL, version, arch)
}

func (s *Service) GetWindowsNpxCmdURL(version, goarch string) string {
	arch := goarch
	if arch == "amd64" {
		arch = "x64"
	}

	return fmt.Sprintf("%s/%s/win-%s/npx.cmd", s.baseURL, version, arch)
}

package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	APP_NAME      = "gnode"
	NODE_DIST_URL = "https://nodejs.org/dist"
)

type Config struct {
	HomeDir    string
	AppDir     string
	CurrentDir string
	GOOS       string
	GOARCH     string
}

func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(homeDir, "."+APP_NAME)
	currentDir := filepath.Join(appDir, "current")

	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "64"
	}

	return &Config{
		HomeDir:    homeDir,
		AppDir:     appDir,
		CurrentDir: currentDir,
		GOOS:       runtime.GOOS,
		GOARCH:     arch,
	}, nil
}

func (c *Config) VersionsDir() string {
	return filepath.Join(c.AppDir, "versions")
}

func (c *Config) GetVersionDir(version string) string {
	return filepath.Join(c.VersionsDir(), version)
}

func (c *Config) GetDistURL() string {
	return NODE_DIST_URL
}

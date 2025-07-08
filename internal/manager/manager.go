package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/joaomarcosfurtado/gnode/internal/downloader"
	"github.com/joaomarcosfurtado/gnode/internal/extractor"
	"github.com/joaomarcosfurtado/gnode/internal/version"
	"github.com/joaomarcosfurtado/gnode/pkg/config"
)

type Manager struct {
	config     *config.Config
	downloader *downloader.Downloader
	extractor  *extractor.Extractor
	version    *version.Service
}

func NewManager(cfg *config.Config) (*Manager, error) {
	return &Manager{
		config:     cfg,
		downloader: downloader.NewDownloader(),
		extractor:  extractor.NewExtractor(),
		version:    version.NewService(cfg.GetDistURL()),
	}, nil
}

func (m *Manager) Init() error {
	if err := os.MkdirAll(m.config.AppDir, 0755); err != nil {
		return fmt.Errorf("Error creating home directory: %v", err)
	}

	if err := os.MkdirAll(m.config.VersionsDir(), 0755); err != nil {
		return fmt.Errorf("Error creating directory versions: %v", err)
	}

	return nil
}

func (m *Manager) Install(versionStr string) error {
	version := m.version.NormalizeVersion(versionStr)
	fmt.Printf("Installing Node.js %s...\n", version)

	versionDir := m.config.GetVersionDir(version)
	if _, err := os.Stat(versionDir); err == nil {
		fmt.Printf("Node.js %s already installed\n", version)
	}

	downloadURL := m.version.GetDownloadURL(version, m.config.GOOS, m.config.GOARCH)

	reader, err := m.downloader.Download(downloadURL)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("error creating version directory: %v", err)
	}

	if err := m.extractor.ExtractTarGz(reader, versionDir); err != nil {
		os.RemoveAll(versionDir)
		return err
	}

	fmt.Printf("Node.js %s installed with success\n", version)
	return nil
}

func (m *Manager) Use(versionStr string) error {
	version := m.version.NormalizeVersion(versionStr)
	versionDir := m.config.GetVersionDir(version)

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		fmt.Errorf("node.js %s is not installed. Execute 'gnode install %v first", version, versionStr)
	}

	if _, err := os.Lstat(m.config.CurrentDir); err != nil {
		return fmt.Errorf("error creating simbolic link %v", err)
	}

	fmt.Printf("Now using Node.js %s\n", version)
	return nil
}

func (m *Manager) ListLocal() error {
	versions, err := m.getLocalVersions()
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		fmt.Println("No version installed")
		return nil
	}

	current, _ := m.getCurrentVersion()
	fmt.Println("Versions installed:")
	for _, v := range versions {
		marker := "  "
		if v == current {
			marker = "* "
		}
		fmt.Printf("%s%s\n", marker, v)
	}

	return nil
}

func (m *Manager) ListRemote() error {
	versions, err := m.version.ListRemote()
	if err != nil {
		return err
	}

	limit := 20
	if len(versions) > limit {
		fmt.Printf("Available versions: (first %d):\n", limit)
		for i := 0; i < limit; i++ {
			fmt.Printf("%s\n", versions[i].Version)
		}
		fmt.Println("... and more")
	} else {
		fmt.Printf("Available versions: (%d):\n", len(versions))
		for _, v := range versions {
			fmt.Printf("%s\n", v.Version)
		}
	}

	return nil
}

func (m *Manager) ShowCurrent() error {
	current, err := m.getCurrentVersion()
	if err != nil {
		return err
	}

	fmt.Println(current)
	return nil
}

func (m *Manager) ShowWhich() error {
	_, err := m.getCurrentVersion()
	if err != nil {
		return err
	}

	nodePath := filepath.Join(m.config.CurrentDir, "bin", "node")
	fmt.Println(nodePath)
	return nil
}

func (m *Manager) Uninstall(versionStr string) error {
	version := m.version.NormalizeVersion(versionStr)
	versionDir := m.config.GetVersionDir(version)

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("node.js %s is not installed", version)
	}

	if current, err := m.getCurrentVersion(); err == nil && current == version {
		return fmt.Errorf("it is not possible to uninstall the current version (%s). use other version first", version)
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return fmt.Errorf("error removing version %v", err)
	}

	fmt.Printf("Node.js uninstalled with success!\n", version)
	return nil
}

func (m *Manager) getLocalVersions() ([]string, error) {
	entries, err := os.ReadDir(m.config.VersionsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error reading directory versions: %v", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	sort.Strings(versions)
	return versions, nil
}

func (m *Manager) getCurrentVersion() (string, error) {
	target, err := os.Readlink(m.config.CurrentDir)
	if err != nil {
		return "", fmt.Errorf("no version of node.js is being used.")
	}

	version := filepath.Base(target)
	return version, nil
}

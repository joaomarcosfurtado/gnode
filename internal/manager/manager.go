package manager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

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

func (m *Manager) Install(versionStr string) error {
	version := m.version.NormalizeVersion(versionStr)
	fmt.Printf("Installing Node.js %s...\n", version)

	versionDir := m.config.GetVersionDir(version)
	if _, err := os.Stat(versionDir); err == nil {
		fmt.Printf("Node.js %s already installed\n", version)
		return nil
	}

	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return fmt.Errorf("error creating version directory: %v", err)
	}

	if runtime.GOOS == "windows" {
		return m.installWindows(version, versionDir)
	}

	downloadURL := m.version.GetDownloadURL(version, m.config.GOOS, m.config.GOARCH)
	reader, err := m.downloader.Download(downloadURL)
	if err != nil {
		os.RemoveAll(versionDir)
		return err
	}
	defer reader.Close()

	if err := m.extractor.ExtractTarGz(reader, versionDir); err != nil {
		os.RemoveAll(versionDir)
		return err
	}

	fmt.Printf("Node.js %s installed with success\n", version)
	return nil
}

func (m *Manager) installWindows(version, versionDir string) error {
	fmt.Printf("Checking available download options...\n")

	strategy := m.version.GetDownloadStrategy(version, m.config.GOARCH)

	switch strategy {
	case "zip":
		fmt.Printf("Using ZIP distribution...\n")
		return m.installWindowsZip(version, versionDir)

	case "binaries":
		fmt.Printf("Using individual binaries...\n")
		return m.installWindowsBinaries(version, versionDir)

	default:
		return fmt.Errorf("no compatible download found for Node.js %s", version)
	}
}

func (m *Manager) installWindowsZip(version, versionDir string) error {
	downloadURL := m.version.GetDownloadURL(version, m.config.GOOS, m.config.GOARCH)
	reader, err := m.downloader.Download(downloadURL)
	if err != nil {
		return fmt.Errorf("error downloading ZIP: %v", err)
	}
	defer reader.Close()

	tempZip := filepath.Join(versionDir, "temp.zip")
	tempFile, err := os.Create(tempZip)
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}

	_, err = io.Copy(tempFile, reader)
	tempFile.Close()
	if err != nil {
		os.Remove(tempZip)
		return fmt.Errorf("error saving ZIP: %v", err)
	}

	fmt.Printf("Extracting ZIP...\n")
	if err := m.extractor.ExtractZip(tempZip, versionDir); err != nil {
		os.Remove(tempZip)
		return fmt.Errorf("error extracting ZIP: %v", err)
	}

	os.Remove(tempZip)

	fmt.Printf("✓ Node.js %s installed successfully from ZIP\n", version)
	return nil
}

func (m *Manager) installWindowsBinaries(version, versionDir string) error {
	fmt.Printf("Downloading Node.js binaries...\n")

	available := m.version.CheckAvailableFiles(version, m.config.GOARCH)

	files := []struct {
		key      string
		getURL   func(string, string) string
		filename string
		required bool
	}{
		{"node", m.version.GetWindowsNodeURL, "node.exe", true},
		{"npm", m.version.GetWindowsNpmURL, "npm", false},
		{"npm.cmd", m.version.GetWindowsNpmCmdURL, "npm.cmd", false},
		{"npx", m.version.GetWindowsNpxURL, "npx", false},
		{"npx.cmd", m.version.GetWindowsNpxCmdURL, "npx.cmd", false},
	}

	downloadedFiles := 0

	for _, file := range files {
		if !available[file.key] {
			if file.required {
				return fmt.Errorf("required file %s not found for Node.js %s", file.filename, version)
			}
			fmt.Printf("Skipping %s (not available)\n", file.filename)
			continue
		}

		url := file.getURL(version, m.config.GOARCH)
		fmt.Printf("Downloading %s...\n", file.filename)

		reader, err := m.downloader.Download(url)
		if err != nil {
			if file.required {
				return fmt.Errorf("error downloading required file %s: %v", file.filename, err)
			}
			fmt.Printf("Warning: failed to download %s: %v\n", file.filename, err)
			continue
		}

		filePath := filepath.Join(versionDir, file.filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			reader.Close()
			return fmt.Errorf("error creating file %s: %v", file.filename, err)
		}

		_, err = io.Copy(outFile, reader)
		outFile.Close()
		reader.Close()

		if err != nil {
			return fmt.Errorf("error writing file %s: %v", file.filename, err)
		}

		downloadedFiles++
		fmt.Printf("✓ Downloaded %s\n", file.filename)
	}

	if err := m.createWindowsWrappers(versionDir, available); err != nil {
		return fmt.Errorf("error creating wrapper scripts: %v", err)
	}

	fmt.Printf("✓ Node.js %s installed successfully (%d files)\n", version, downloadedFiles)

	if !available["npm"] && !available["npm.cmd"] {
		fmt.Printf("⚠️  npm not available for this version\n")
		fmt.Printf("   You can install it manually: npm install -g npm\n")
	}

	return nil
}

func (m *Manager) createWindowsWrappers(versionDir string, available map[string]bool) error {
	if available["npm"] && !available["npm.cmd"] {
		npmBat := filepath.Join(versionDir, "npm.bat")
		npmContent := `@echo off
setlocal
set "NODE_EXE=%~dp0node.exe"
set "NPM_CLI_JS=%~dp0npm"
if exist "%NPM_CLI_JS%" (
  "%NODE_EXE%" "%NPM_CLI_JS%" %*
) else (
  echo npm not found in this Node.js installation
  exit /b 1
)
`
		if err := os.WriteFile(npmBat, []byte(npmContent), 0644); err != nil {
			return fmt.Errorf("error creating npm.bat: %v", err)
		}
		fmt.Printf("✓ Created npm.bat wrapper\n")
	}

	if available["npx"] && !available["npx.cmd"] {
		npxBat := filepath.Join(versionDir, "npx.bat")
		npxContent := `@echo off
setlocal
set "NODE_EXE=%~dp0node.exe"
set "NPX_CLI_JS=%~dp0npx"
if exist "%NPX_CLI_JS%" (
  "%NODE_EXE%" "%NPX_CLI_JS%" %*
) else (
  echo npx not found in this Node.js installation
  exit /b 1
)
`
		if err := os.WriteFile(npxBat, []byte(npxContent), 0644); err != nil {
			return fmt.Errorf("error creating npx.bat: %v", err)
		}
		fmt.Printf("✓ Created npx.bat wrapper\n")
	}

	return nil
}

func (m *Manager) Use(versionStr string, printEnv bool) error {
	version := m.version.NormalizeVersion(versionStr)
	versionDir := m.config.GetVersionDir(version)

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return fmt.Errorf("node.js %s is not installed. Execute 'gnode install %v' first", version, versionStr)
	}

	if err := m.ensureInSystemPath(); err != nil {
		fmt.Printf("Warning: could not add to PATH: %v\n", err)
	}

	if err := m.updateCurrentVersion(version); err != nil {
		return fmt.Errorf("error updating current version: %v", err)
	}

	fmt.Printf("Now using Node.js %s\n", version)

	if m.needsPathRefresh() {
		fmt.Printf("Please restart your terminal or run: refreshenv\n")
	}

	return nil
}

func (m *Manager) ensureInSystemPath() error {
	if runtime.GOOS != "windows" {
		return m.ensureInUnixPath()
	}

	currentDir := m.config.CurrentDir

	if m.isInPath(currentDir) {
		return nil
	}

	fmt.Printf("Adding gnode to system PATH...\n")

	return m.addToWindowsPath(currentDir)
}

func (m *Manager) isInPath(dir string) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("reg", "query", "HKCU\\Environment", "/v", "PATH")
		output, err := cmd.Output()
		if err != nil {
			return false
		}
		return strings.Contains(string(output), dir)
	}

	path := os.Getenv("PATH")
	return strings.Contains(path, dir)
}

func (m *Manager) addToWindowsPath(dir string) error {
	cmd := exec.Command("reg", "query", "HKCU\\Environment", "/v", "PATH")
	output, err := cmd.Output()

	currentPath := ""
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "PATH") && strings.Contains(line, "REG_") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					currentPath = strings.Join(parts[2:], " ")
					break
				}
			}
		}
	}

	newPath := currentPath
	if newPath != "" {
		newPath += ";" + dir
	} else {
		newPath = dir
	}

	cmd = exec.Command("reg", "add", "HKCU\\Environment", "/v", "PATH", "/t", "REG_EXPAND_SZ", "/d", newPath, "/f")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update PATH: %v", err)
	}

	m.broadcastPathChange()

	return nil
}

func (m *Manager) broadcastPathChange() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command",
			"[System.Environment]::SetEnvironmentVariable('PATH', [System.Environment]::GetEnvironmentVariable('PATH', 'User'), 'User')")
		cmd.Run()
	}
}

func (m *Manager) needsPathRefresh() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	currentDir := m.config.CurrentDir
	currentPath := os.Getenv("PATH")

	return !strings.Contains(currentPath, currentDir)
}

func (m *Manager) ensureInUnixPath() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	currentDir := m.config.CurrentDir

	shell := os.Getenv("SHELL")
	var rcFile string

	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(homeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(homeDir, ".bashrc")
	}

	line := fmt.Sprintf(`export PATH="%s:$PATH"`, currentDir)

	content, err := os.ReadFile(rcFile)
	if err == nil && strings.Contains(string(content), line) {
		return nil
	}

	return appendLineIfNotExists(rcFile, line)
}

func (m *Manager) updateCurrentVersion(version string) error {
	currentPath := m.config.CurrentDir
	versionPath := m.config.GetVersionDir(version)

	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return fmt.Errorf("version directory does not exist: %s", versionPath)
	}

	if _, err := os.Lstat(currentPath); err == nil {
		os.RemoveAll(currentPath)
		if runtime.GOOS == "windows" {
			exec.Command("cmd", "/c", "rmdir", "/Q", currentPath).Run()
		}
	}

	if runtime.GOOS == "windows" {
		return m.createWindowsJunction(versionPath, currentPath)
	}

	return os.Symlink(versionPath, currentPath)
}

func (m *Manager) createWindowsJunction(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("cannot create parent directory: %v", err)
	}

	cmd := exec.Command("cmd", "/c", "mklink", "/J", fmt.Sprintf(`"%s"`, dst), fmt.Sprintf(`"%s"`, src))
	if err := cmd.Run(); err == nil {
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
	}

	psCmd := fmt.Sprintf(`New-Item -ItemType Junction -Path "%s" -Target "%s" -Force`, dst, src)
	cmd = exec.Command("powershell", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		if _, err := os.Stat(dst); err == nil {
			return nil
		}
	}

	fmt.Printf("Warning: Could not create junction, copying directory instead\n")
	return copyDir(src, dst)
}

func (m *Manager) ShowCurrent() error {
	currentDir := m.config.CurrentDir

	if _, err := os.Stat(currentDir); os.IsNotExist(err) {
		return fmt.Errorf("no version of node.js is being used")
	}

	var nodeExe string
	if runtime.GOOS == "windows" {
		nodeExe = filepath.Join(currentDir, "node.exe")
	} else {
		nodeExe = filepath.Join(currentDir, "node")
	}

	if _, err := os.Stat(nodeExe); os.IsNotExist(err) {
		return fmt.Errorf("no version of node.js is being used")
	}

	cmd := exec.Command(nodeExe, "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no version of node.js is being used")
	}

	version := strings.TrimSpace(string(output))
	fmt.Println(version)
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

func (m *Manager) ShowWhich() error {
	_, err := m.getCurrentVersion()
	if err != nil {
		return err
	}

	var nodePath string
	if runtime.GOOS == "windows" {
		nodePath = filepath.Join(m.config.CurrentDir, "node.exe")
	} else {
		nodePath = filepath.Join(m.config.CurrentDir, "bin", "node")
	}

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

	fmt.Printf("Node.js %s uninstalled with success!\n", version)
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

func (m *Manager) getCurrentVersionWindows() (string, error) {
	currentDir := m.config.CurrentDir

	nodeExe := filepath.Join(currentDir, "node.exe")
	if _, err := os.Stat(nodeExe); err == nil {
		cmd := exec.Command(nodeExe, "--version")
		output, err := cmd.Output()
		if err == nil {
			version := strings.TrimSpace(string(output))
			return version, nil
		}
	}

	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("try { (Get-Item '%s' -Force).Target } catch { $null }", currentDir))
	output, err := cmd.Output()
	if err == nil {
		target := strings.TrimSpace(string(output))
		if target != "" && target != "null" && target != currentDir {
			return filepath.Base(target), nil
		}
	}

	return "", fmt.Errorf("no version of node.js is being used")
}

func (m *Manager) getCurrentVersion() (string, error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("node", "--version")
		output, err := cmd.Output()
		if err == nil {
			version := strings.TrimSpace(string(output))
			return version, nil
		}

		return m.getCurrentVersionWindows()
	}

	target, err := os.Readlink(m.config.CurrentDir)
	if err != nil {
		return "", fmt.Errorf("no version of node.js is being used")
	}

	version := filepath.Base(target)
	return version, nil
}

func (m *Manager) Status() error {
	fmt.Printf("gnode status:\n")

	currentDir := m.config.CurrentDir
	if m.isInPath(currentDir) {
		fmt.Printf("✓ gnode is in system PATH\n")
	} else {
		fmt.Printf("✗ gnode is NOT in system PATH\n")
		fmt.Printf("  Run: gnode use <version> to add automatically\n")
	}

	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err == nil {
		nodeVersion := strings.TrimSpace(string(output))
		fmt.Printf("✓ Active version: %s\n", nodeVersion)
		fmt.Printf("✓ Node.js accessible: %s\n", nodeVersion)

		npmCmd := exec.Command("npm", "--version")
		npmOutput, npmErr := npmCmd.Output()
		if npmErr == nil {
			npmVersion := strings.TrimSpace(string(npmOutput))
			fmt.Printf("✓ npm accessible: %s\n", npmVersion)
		} else {
			fmt.Printf("⚠️  npm not accessible\n")
		}
	} else {
		fmt.Printf("✗ No active version\n")
		fmt.Printf("✗ Node.js not accessible\n")
		fmt.Printf("  Try: gnode use <version>\n")
	}

	return nil
}

func appendLineIfNotExists(filePath, line string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == line {
			return nil
		}
	}

	if _, err := file.WriteString("\n" + line + "\n"); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Init() error {
	if err := os.MkdirAll(m.config.AppDir, 0755); err != nil {
		return fmt.Errorf("error creating home directory: %v", err)
	}

	if err := os.MkdirAll(m.config.VersionsDir(), 0755); err != nil {
		return fmt.Errorf("error creating directory versions: %v", err)
	}

	emptyPath := filepath.Join(m.config.AppDir, "empty")
	currentPath := m.config.CurrentDir

	if _, err := os.Stat(emptyPath); os.IsNotExist(err) {
		if err := os.Mkdir(emptyPath, 0755); err != nil {
			return fmt.Errorf("error creating empty dir: %v", err)
		}
	}

	os.RemoveAll(currentPath)

	if runtime.GOOS == "windows" {
		err := createJunction(emptyPath, currentPath)
		if err != nil {
			errCopy := copyDir(emptyPath, currentPath)
			if errCopy != nil {
				return fmt.Errorf("error creating junction or copying dir: %v, %v", err, errCopy)
			}
			return nil
		}
		return nil
	}

	if err := os.Symlink(emptyPath, currentPath); err != nil {
		return fmt.Errorf("error creating symlink: %v", err)
	}

	return nil
}

func createJunction(src, dst string) error {
	cmd := exec.Command("cmd", "/c", "mklink", "/J", dst, src)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mklink /J error: %v, output: %s", err, string(out))
	}
	return nil
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

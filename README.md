# gnode - Node.js Version Manager

A fast and simple Node.js version manager for Windows, macOS, and Linux. Works like nvm-windows with automatic PATH management.

## Features

- 🚀 **Fast and lightweight** - Single executable, no dependencies
- 🎯 **Cross-platform** - Windows, macOS, and Linux
- 📋 **npm included** - Automatically downloads npm when available

## Installation

### Quick Install (Windows)

1. Download the latest release from GitHub
2. Extract `gnode.exe` and `install.bat`
3. Run `install.bat` (adds to PATH automatically)
4. Restart your terminal

### Manual Install

#### Windows
```bash
# Download gnode.exe to a folder in your PATH
# Or use the install.bat script
```

#### macOS/Linux
```bash
# Download the binary
chmod +x gnode
sudo mv gnode /usr/local/bin/
```

## Usage

### First Time Setup
```bash
# Install a Node.js version
gnode install v20.12.0

# Use it (automatically configures PATH on first use)
gnode use v20.12.0

# Restart terminal if needed (Windows only)
node -v  # Should show v20.12.0
npm -v   # npm is included
```

### Daily Usage
```bash
# Install multiple versions
gnode install v18.19.1
gnode install v20.12.0
gnode install v22.0.0

# Switch between versions instantly
gnode use v18.19.1
node -v  # v18.19.1

gnode use v22.0.0  
node -v  # v22.0.0

# List installed versions
gnode list
# Output:
#   v18.19.1
# * v20.12.0
#   v22.0.0

# List available versions to install
gnode list-remote

# Check system status
gnode status
# Output:
# ✓ gnode is in system PATH
# ✓ Active version: v20.12.0
# ✓ Node.js accessible: v20.12.0

# Remove a version
gnode uninstall v18.19.1
```

## Commands

| Command | Description |
|---------|-------------|
| `gnode install <version>` | Install Node.js version |
| `gnode use <version>` | Switch to Node.js version |
| `gnode list` | List installed versions |
| `gnode list-remote` | List available versions |
| `gnode current` | Show current version |
| `gnode which` | Show Node.js executable path |
| `gnode uninstall <version>` | Remove Node.js version |
| `gnode status` | Show gnode status |
| `gnode help` | Show help |

## How it Works

gnode works similarly to nvm-windows:

1. **Installation**: Downloads Node.js to `~/.gnode/versions/`
2. **PATH Management**: Adds `~/.gnode/current` to system PATH (once)
3. **Version Switching**: Uses symlinks/junctions to point `current` to active version
4. **Instant Switching**: Changing versions just updates the symlink

```
~/.gnode/
├── current/          # Symlink to active version
├── versions/
│   ├── v18.19.1/
│   ├── v20.12.0/
│   └── v22.0.0/
└── ...
```

## Building from Source

### Prerequisites
- Go 1.19 or later
- Git

### Build for Current Platform
```bash
git clone https://github.com/joaomarcosfurtado/gnode.git
cd gnode
go build -o gnode ./cmd/gnode
```

### Cross-Compilation

#### Build for Windows (from macOS/Linux)
```bash
GOOS=windows GOARCH=amd64 go build -o dist/windows/gnode.exe ./cmd/gnode
```

#### Build for macOS (from Windows/Linux)
```bash
GOOS=darwin GOARCH=amd64 go build -o dist/macos/gnode ./cmd/gnode
```

#### Build for Linux (from Windows/macOS)
```bash
GOOS=linux GOARCH=amd64 go build -o dist/linux/gnode ./cmd/gnode
```

#### Build for All Platforms
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o dist/windows/gnode.exe ./cmd/gnode

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o dist/macos/gnode-intel ./cmd/gnode

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o dist/macos/gnode-arm64 ./cmd/gnode

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/linux/gnode ./cmd/gnode
```

## Project Structure

```
gnode/
├── cmd/gnode/           # Main application entry point
├── internal/
│   ├── downloader/      # HTTP download functionality
│   ├── extractor/       # Archive extraction (tar.gz, zip)
│   ├── manager/         # Core version management logic
│   └── version/         # Version service and URL handling
├── pkg/config/          # Configuration management
├── dist/                # Build outputs
│   ├── windows/
│   ├── macos/
│   └── linux/
├── install.bat          # Windows installation script
└── README.md
```

## Troubleshooting

### Windows

**Node command not found after installation:**
```bash
# Restart terminal or run:
refreshenv

# Check status:
gnode status

# Manually add to PATH if needed:
gnode use v20.12.0
```

**Permission errors:**
```bash
# Run install.bat as administrator, or
# Manually copy gnode.exe to C:\Windows\System32\
```

### macOS/Linux

**Permission denied:**
```bash
chmod +x gnode
```

**Command not found:**
```bash
# Make sure gnode is in your PATH:
echo $PATH

# Add to shell profile if needed:
echo 'export PATH="$HOME/.gnode/current:$PATH"' >> ~/.bashrc
```

### General

**Version not found:**
```bash
# Check available versions:
gnode list-remote

# Use exact version format:
gnode install v20.12.0  # ✓ Correct
gnode install 20.12.0   # ✓ Also works
gnode install 20.12     # ✗ Won't work
```

**npm not available:**
```bash
# Some very new versions may not have npm in binary distribution
# Use LTS versions for guaranteed npm support:
gnode install v20.12.0  # LTS with npm
gnode install v18.19.1  # LTS with npm
```

## Comparison with Other Tools

| Feature         | gnode | nvm-windows | nvm | fnm |
|-----------------|-------|-------------|-----|-----|
| Single binary   | ✅    | ❌          | ❌  | ✅  |
| Auto PATH       | ✅    | ✅          | ❌  | ❌  |
| Cross-platform  | ✅    | ❌          | ✅  | ✅  |
| Fast switching  | ✅    | ✅          | ✅  | ✅  |
| npm included    | ✅    | ✅          | ✅  | ❌  |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [nvm](https://github.com/nvm-sh/nvm) and [nvm-windows](https://github.com/coreybutler/nvm-windows)
- Built with using Go

## Support

- 🐛 **Bug reports**: [GitHub Issues](https://github.com/yourusername/gnode/issues)
- 💬 **Questions**: [GitHub Discussions](https://github.com/yourusername/gnode/discussions)
- 📖 **Documentation**: This README and `gnode help`

---

**Happy coding with gnode! 🚀**
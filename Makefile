BINARY_NAME=gnode
INSTALL_PATH=$(HOME)/.gnode
SHELL_SCRIPT=gnode.sh

.PHONY: all build install install-auto clean uninstall test help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd/gnode

install: build
	@echo "Installing $(BINARY_NAME)..."
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	cp scripts/$(SHELL_SCRIPT) $(INSTALL_PATH)/$(SHELL_SCRIPT)
	chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	chmod +x $(INSTALL_PATH)/$(SHELL_SCRIPT)
	@echo ""
	@echo "$(BINARY_NAME) installed successfully!"
	@echo ""
	@echo "To use GNode, add the following line to your ~/.bashrc or ~/.zshrc:"
	@echo "source $(INSTALL_PATH)/$(SHELL_SCRIPT)"
	@echo ""
	@echo "Or run:"
	@echo "echo 'source $(INSTALL_PATH)/$(SHELL_SCRIPT)' >> ~/.bashrc"
	@echo "source ~/.bashrc"

install-auto: install
	@echo "Configuring shell automatically..."
	@if [ -f ~/.bashrc ]; then \
		echo "source $(INSTALL_PATH)/$(SHELL_SCRIPT)" >> ~/.bashrc; \
		echo "Configuration added to ~/.bashrc"; \
	elif [ -f ~/.zshrc ]; then \
		echo "source $(INSTALL_PATH)/$(SHELL_SCRIPT)" >> ~/.zshrc; \
		echo "Configuration added to ~/.zshrc"; \
	else \
		echo "Shell configuration file not found"; \
		echo "Add manually: source $(INSTALL_PATH)/$(SHELL_SCRIPT)"; \
	fi
	@echo ""
	@echo "Restart your terminal or run: source ~/.bashrc (or ~/.zshrc)"

clean:
	@echo "Cleaning compiled files..."
	rm -f $(BINARY_NAME)

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -rf $(INSTALL_PATH)
	@echo "$(BINARY_NAME) uninstalled successfully!"
	@echo "Manually remove the line 'source $(INSTALL_PATH)/$(SHELL_SCRIPT)' from your ~/.bashrc or ~/.zshrc"

test:
	@echo "Running tests..."
	go test ./...

help:
	@echo "Available commands:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install GNode"
	@echo "  make install-auto - Install and configure automatically"
	@echo "  make clean       - Remove compiled files"
	@echo "  make uninstall   - Remove GNode"
	@echo "  make test        - Run tests"
	@echo "  make help        - Show this help"
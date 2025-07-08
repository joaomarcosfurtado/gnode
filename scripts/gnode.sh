#!/bin/bash

# GNode (Go Node Version Manager) Shell Integration
export GNODE_DIR="$HOME/.gnode"
export GNODE_BIN="$GNODE_DIR/gnode"

# Main gnode function
gnode() {
    case "$1" in
        "use")
            # Execute the use command
            $GNODE_BIN use "$2"
            if [ $? -eq 0 ]; then
                # Update PATH if command was successful
                # Remove old node paths to avoid duplication
                export PATH=$(echo "$PATH" | sed -e "s|$GNODE_DIR/versions/[^/]*/bin:||g")
                export PATH="$GNODE_DIR/current/bin:$PATH"
                
                # Check if node is working
                if command -v node >/dev/null 2>&1; then
                    echo "PATH updated. Node.js $(node --version) is now active."
                fi
            fi
            ;;
        *)
            # For all other commands, just execute the binary
            $GNODE_BIN "$@"
            ;;
    esac
}

# Auto-initialization - add current version to PATH if it exists
if [ -L "$GNODE_DIR/current" ]; then
    export PATH="$GNODE_DIR/current/bin:$PATH"
fi

# Bash completion
if [ -n "$BASH_VERSION" ]; then
    _gnode_completion() {
        local cur prev opts
        COMPREPLY=()
        cur="${COMP_WORDS[COMP_CWORD]}"
        prev="${COMP_WORDS[COMP_CWORD-1]}"
        opts="install use list list-remote current which uninstall help"

        if [ $COMP_CWORD -eq 1 ]; then
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        elif [ $COMP_CWORD -eq 2 ]; then
            case "${prev}" in
                use|uninstall)
                    local versions=$(gnode list 2>/dev/null | grep -E "^[[:space:]]*\*?[[:space:]]*v" | sed 's/^[[:space:]]*\*\?[[:space:]]*//')
                    COMPREPLY=( $(compgen -W "${versions}" -- ${cur}) )
                    ;;
                install)
                    local common_versions="18.17.0 18.18.0 20.5.0 20.6.0 20.7.0"
                    COMPREPLY=( $(compgen -W "${common_versions}" -- ${cur}) )
                    ;;
            esac
        fi
    }
    complete -F _gnode_completion gnode
fi

# Zsh completion
if [ -n "$ZSH_VERSION" ]; then
    _gnode_zsh_completion() {
        local -a commands
        commands=(
            'install:Install a Node.js version'
            'use:Use a Node.js version'
            'list:List installed versions'
            'list-remote:List available versions for download'
            'current:Show current version'
            'which:Show Node.js executable path'
            'uninstall:Uninstall a Node.js version'
            'help:Show help'
        )
        
        if (( CURRENT == 2 )); then
            _describe 'commands' commands
        elif (( CURRENT == 3 )); then
            case "$words[2]" in
                use|uninstall)
                    local versions
                    versions=($(gnode list 2>/dev/null | grep -E "^[[:space:]]*\*?[[:space:]]*v" | sed 's/^[[:space:]]*\*\?[[:space:]]*//' | tr '\n' ' '))
                    _describe 'versions' versions
                    ;;
                install)
                    local common_versions
                    common_versions=("18.17.0" "18.18.0" "20.5.0" "20.6.0" "20.7.0")
                    _describe 'common versions' common_versions
                    ;;
            esac
        fi
    }
    
    compdef _gnode_zsh_completion gnode
fi
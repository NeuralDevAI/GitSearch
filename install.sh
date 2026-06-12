#!/bin/bash
# install.sh
# Installs gitsearch globally
echo -e "\033[36mInstalling gitsearch...\033[0m"

# Check if go is installed
if command -v go &> /dev/null; then
    echo -e "\033[32mGo compiler found. Compiling and installing...\033[0m"
    go install .
    if [ $? -eq 0 ]; then
        echo -e "\033[32mgitsearch has been successfully installed via 'go install'!\033[0m"
        echo -e "\033[33mMake sure \$(go env GOPATH)/bin is in your PATH.\033[0m"
        exit 0
    else
        echo -e "\033[31mFailed to install gitsearch using 'go install'.\033[0m" >&2
        exit 1
    fi
else
    echo -e "\033[33mGo compiler not found. Trying to install precompiled binary...\033[0m"
    # Check if precompiled binary exists
    if [ -f "./gitsearch" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        cp ./gitsearch "$INSTALL_DIR/gitsearch"
        chmod +x "$INSTALL_DIR/gitsearch"
        echo -e "\033[32mCopied gitsearch to $INSTALL_DIR\033[0m"
        
        # Check PATH
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            echo -e "\033[33mPlease add $INSTALL_DIR to your PATH. For example, add this to your .bashrc or .zshrc:\033[0m"
            echo "export PATH=\"\$PATH:$INSTALL_DIR\""
        fi
        echo -e "\033[32mgitsearch has been successfully installed!\033[0m"
        exit 0
    else
        echo -e "\033[31mError: Neither Go compiler nor precompiled 'gitsearch' was found in the current directory.\033[0m" >&2
        exit 1
    fi
fi

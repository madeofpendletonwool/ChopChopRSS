#!/bin/bash
set -e

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Convert GitHub repo format if provided
REPO=${1:-"madeofpendletonwool/chopchoprss"}
REPO_OWNER=$(echo "$REPO" | cut -d "/" -f1)
REPO_NAME=$(echo "$REPO" | cut -d "/" -f2)

# Fetch latest version
VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Error: Could not determine latest version. Please check the repository name and your internet connection."
  exit 1
fi

echo "Installing ChopChopRSS $VERSION for $OS/$ARCH..."

# Download binary
TMP_DIR=$(mktemp -d)
BINARY_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/${VERSION}/chopchoprss_${OS}_${ARCH}.tar.gz"
echo "Downloading from $BINARY_URL"

curl -sL "$BINARY_URL" -o "${TMP_DIR}/chopchoprss.tar.gz"
tar -xzf "${TMP_DIR}/chopchoprss.tar.gz" -C "${TMP_DIR}"

# Install binary
INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"
mv "${TMP_DIR}/chopchoprss" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/chopchoprss"

# Setup shell completion if shell detected
if [ -n "$SHELL" ]; then
  SHELL_NAME=$(basename "$SHELL")

  if [ "$SHELL_NAME" = "bash" ]; then
    echo "Setting up Bash completion..."
    COMPLETION_DIR="$HOME/.bash_completion.d"
    mkdir -p "$COMPLETION_DIR"
    "$INSTALL_DIR/chopchoprss" completion bash > "$COMPLETION_DIR/chopchoprss"

    # Add to .bashrc if not already there
    if ! grep -q "bash_completion.d/chopchoprss" "$HOME/.bashrc"; then
      echo "[ -f $COMPLETION_DIR/chopchoprss ] && source $COMPLETION_DIR/chopchoprss" >> "$HOME/.bashrc"
      echo "Added completion to .bashrc. Please restart your shell or run: source ~/.bashrc"
    fi

  elif [ "$SHELL_NAME" = "zsh" ]; then
    echo "Setting up Zsh completion..."
    ZSH_COMPLETION_DIR="$HOME/.zsh/completion"
    mkdir -p "$ZSH_COMPLETION_DIR"
    "$INSTALL_DIR/chopchoprss" completion zsh > "$ZSH_COMPLETION_DIR/_chopchoprss"

    # Add to .zshrc if not already there
    if ! grep -q "fpath=($ZSH_COMPLETION_DIR" "$HOME/.zshrc"; then
      echo "fpath=($ZSH_COMPLETION_DIR \$fpath)" >> "$HOME/.zshrc"
      echo "autoload -U compinit; compinit" >> "$HOME/.zshrc"
      echo "Added completion to .zshrc. Please restart your shell or run: source ~/.zshrc"
    fi

  elif [ "$SHELL_NAME" = "fish" ]; then
    echo "Setting up Fish completion..."
    FISH_COMPLETION_DIR="$HOME/.config/fish/completions"
    mkdir -p "$FISH_COMPLETION_DIR"
    "$INSTALL_DIR/chopchoprss" completion fish > "$FISH_COMPLETION_DIR/chopchoprss.fish"
    echo "Fish completion installed. It will be active in new shell sessions."
  fi
fi

# Check if the installation directory is in PATH
if ! echo "$PATH" | tr ':' '\n' | grep -q "$INSTALL_DIR"; then
  echo
  echo "⚠️  Warning: $INSTALL_DIR is not in your PATH."
  echo "Please add it to your PATH by adding this line to your shell configuration file:"
  echo
  echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
  echo
else
  echo
  echo "🎉 ChopChopRSS installed successfully in $INSTALL_DIR"
  echo "Run 'chopchoprss' to get started"
  echo
fi

# Clean up
rm -rf "$TMP_DIR"

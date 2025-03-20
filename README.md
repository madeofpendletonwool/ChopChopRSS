# ChopChopRSS

<p align="center">
  <img src="chopchop.png" alt="ChopChopRSS Logo" width="200">
</p>

A fast and simple CLI tool for generating and managing RSS feeds.

## Features

- Create and manage multiple RSS feeds
- Quickly add new entries to feeds
- Support for images in feed entries
- Customizable feed metadata
- Serves feeds via HTTP
- Dockerized for easy deployment

## Installation

### Quick Install (Linux/macOS)

```bash
# Install latest version
curl -fsSL https://raw.githubusercontent.com/madeofpendletonwool/chopchoprss/main/install.sh | bash
```

The installer will:
- Download the appropriate binary for your OS and architecture
- Install it to `~/.local/bin/` (ensure this is in your PATH)
- Set up shell completion for your detected shell (bash, zsh, or fish)

### Homebrew (macOS/Linux)

```bash
brew tap madeofpendletonwool/chopchoprss
brew install chopchoprss
```

### From Source

```bash
git clone https://github.com/madeofpendletonwool/chopchoprss.git
cd chopchoprss
go build -o chopchoprss
```

### Using Docker

```bash
# Pull the latest image from GitHub Container Registry
docker pull ghcr.io/madeofpendletonwool/chopchoprss:latest

# Or build it locally
docker build -t chopchoprss .

# Run the server (default behavior)
docker run -p 8090:8090 -v chopchoprss-data:/data ghcr.io/madeofpendletonwool/chopchoprss:latest

# Create a feed
docker run -v chopchoprss-data:/data ghcr.io/madeofpendletonwool/chopchoprss:latest create-feed -n blog -t "My Blog" -d "My awesome blog"

# Add an entry
docker run -v chopchoprss-data:/data ghcr.io/madeofpendletonwool/chopchoprss:latest create-entry -f blog -t "First Post" -c "Hello world!"

# List feeds
docker run -v chopchoprss-data:/data ghcr.io/madeofpendletonwool/chopchoprss:latest list-feeds

# For convenience, you can create an alias
alias chopchoprss-docker="docker run -v chopchoprss-data:/data ghcr.io/madeofpendletonwool/chopchoprss:latest"

# Then use it like the normal CLI
chopchoprss-docker create-entry -f blog -t "Second Post" -c "Another post!"
```

#### Using Docker Compose

```bash
# Start the server
docker-compose up -d

# Run commands
docker-compose exec chopchoprss create-feed -n blog -t "My Blog"
docker-compose exec chopchoprss create-entry -f blog -t "First Post" -c "Content here"
docker-compose exec chopchoprss list-feeds
```

## Installation

### From Source

```bash
git clone https://github.com/madeofpendletonwool/chopchoprss.git
cd chopchoprss
go build -o chopchoprss
```

### Using Docker

```bash
# Build the Docker image
docker build -t chopchoprss .

# Run the container
docker run -p 8090:8090 chopchoprss
```

## Shell Completion

ChopChopRSS provides shell completion for Bash, Zsh, Fish, and PowerShell.

### Bash

```bash
# Generate the completion script
source <(chopchoprss completion bash)

# To load completions for each session, execute once:
# Linux:
chopchoprss completion bash > /etc/bash_completion.d/chopchoprss
# macOS:
chopchoprss completion bash > $(brew --prefix)/etc/bash_completion.d/chopchoprss
```

### Zsh

```bash
# If shell completion is not already enabled, you need to enable it:
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Generate and install the completion script
chopchoprss completion zsh > "${fpath[1]}/_chopchoprss"

# Start a new shell for this to take effect
```

### Fish

```bash
chopchoprss completion fish > ~/.config/fish/completions/chopchoprss.fish
```

### PowerShell

```powershell
# For current session
chopchoprss completion powershell | Out-String | Invoke-Expression

# For all sessions (add to profile)
chopchoprss completion powershell > chopchoprss.ps1
```

## Usage

### Creating a Feed

```bash
chopchoprss create-feed -n feedname -t "Feed Title" -d "Feed Description" -l "https://example.com" -a "Author Name" -e "author@example.com"
```

Options:
- `-n, --name`: Feed name (required, used in URL)
- `-t, --title`: Feed title (required)
- `-d, --description`: Feed description
- `-l, --link`: Feed link
- `-a, --author`: Feed author name
- `-e, --email`: Feed author email

### Adding an Entry to a Feed

```bash
chopchoprss create-entry -f feedname -t "Entry Title" -c "Entry content goes here" -l "https://example.com/entry" -i "https://example.com/image.jpg"
```

Options:
- `-f, --feed`: Feed name (required)
- `-t, --title`: Entry title (required)
- `-c, --content`: Entry content (required)
- `-l, --link`: Entry link
- `-i, --image`: Entry image URL

### Listing All Feeds

```bash
chopchoprss list-feeds
```

### Listing Entries in a Feed

```bash
chopchoprss list-entries -f feedname
```

Options:
- `-f, --feed`: Feed name (required)

### Deleting a Feed

```bash
chopchoprss delete-feed -n feedname
```

Options:
- `-n, --name`: Feed name (required)

### Deleting an Entry

```bash
chopchoprss delete-entry -f feedname -i 0
```

Options:
- `-f, --feed`: Feed name (required)
- `-i, --index`: Entry index (required, zero-based)

### Starting the Server

```bash
chopchoprss serve -p 8090
```

Options:
- `-p, --port`: Server port (default: 8090)

## Accessing Feeds

Once the server is running, feeds are available at:
```
http://localhost:8090/feedname
```

Replace `feedname` with the name you gave to your feed.

## Configuration

ChopChopRSS stores its configuration in the `~/.chopchoprss/config.json` file. This file is created automatically when you first run the application.

## Docker Usage

Build and run the Docker container:

```bash
# Build the image
docker build -t chopchoprss .

# Run the container with a volume to persist data
docker run -p 8090:8090 -v ~/.chopchoprss:/root/.chopchoprss chopchoprss
```

## Examples

### Creating a tech news feed and adding an entry

```bash
# Create a feed
chopchoprss create-feed -n tech -t "Tech News" -d "Latest technology news" -a "Tech Editor" -e "editor@technews.com"

# Add an entry
chopchoprss create-entry -f tech -t "New Golang Release" -c "Go 1.18 has been released with exciting new features." -l "https://example.com/golang-release" -i "https://example.com/golang-logo.png"

# Start the server
chopchoprss serve

# Access the feed at http://localhost:8090/tech
```

## License

GPL3

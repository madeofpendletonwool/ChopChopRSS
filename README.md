# ChopChopRSS

<p align="center">
  <img src="./chopchop.png" alt="ChopChopRSS Logo" width="200">
</p>

A fast and simple CLI tool for generating and managing RSS feeds and podcast feeds from audio directories. This is a great little tool for deploying your own podcast feeds.

## Features

### RSS Feeds
- Create and manage multiple RSS feeds
- Quickly add new entries to feeds
- Support for images in feed entries
- Customizable feed metadata
- Serves feeds via HTTP

### Podcast Feeds
- **Automatic podcast generation** from audio directories
- **Audio metadata extraction** from MP3, M4A, WAV, FLAC, and OGG files
- **iTunes-compatible RSS** with podcast extensions
- **Episode management** with automatic file discovery
- **Multi-podcast support** with isolated feeds
- **Docker-ready** for hosting multiple podcast archives

### General
- HTTP server for feed and audio file serving
- Docker and Docker Compose support
- Shell completion (bash, zsh, fish, powershell)
- Persistent configuration and state

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

ChopChopRSS supports two main use cases: **RSS feeds** for traditional content syndication and **podcast feeds** for serving audio content with proper podcast metadata.

## RSS Feeds

### Creating and Managing RSS Feeds

```bash
# Create a basic RSS feed
chopchoprss create-feed -n blog -t "My Blog" -d "Personal thoughts and updates"

# Create a comprehensive feed with all metadata
chopchoprss create-feed \
  -n tech-news \
  -t "Tech News Daily" \
  -d "Latest technology news and updates" \
  -l "https://example.com" \
  -a "John Doe" \
  -e "john@example.com"

# Add entries to feeds
chopchoprss create-entry \
  -f tech-news \
  -t "New Golang Release" \
  -c "Go 1.21 has been released with exciting new features including improved performance and new standard library additions." \
  -l "https://example.com/golang-release" \
  -i "https://example.com/golang-logo.png"

# List all feeds and their entries
chopchoprss list-feeds
chopchoprss list-entries -f tech-news

# Delete entries or entire feeds
chopchoprss delete-entry -f tech-news -i 0
chopchoprss delete-feed -n old-feed
```

## Podcast Feeds

### Creating Podcasts from Audio Directories

ChopChopRSS can automatically generate podcast feeds by scanning directories containing audio files and extracting metadata.

```bash
# Create a podcast from an audio directory
chopchoprss create-podcast \
  -n "my-podcast" \
  -t "My Amazing Podcast" \
  -d "Weekly discussions about technology, life, and everything in between" \
  -a "Host Name" \
  -e "host@example.com" \
  -u "http://localhost:8090/my-podcast" \
  -r "/path/to/audio/episodes" \
  -c "Technology" \
  -i "https://example.com/podcast-cover.jpg" \
  --language "en" \
  --copyright "© 2024 My Podcast"

# Refresh podcast when new episodes are added
chopchoprss refresh-podcast -n "my-podcast"

# List all podcasts
chopchoprss list-podcasts

# Delete a podcast
chopchoprss delete-podcast -n "old-podcast"
```

#### Supported Audio Formats
- **MP3** (.mp3) - Most common podcast format
- **M4A** (.m4a) - Apple's AAC format
- **WAV** (.wav) - Uncompressed audio
- **FLAC** (.flac) - Lossless compression
- **OGG** (.ogg) - Open source format

#### Automatic Metadata Extraction
ChopChopRSS automatically extracts episode information from audio file metadata:
- **Episode title** from ID3 title tag (or filename if missing)
- **Episode description** from ID3 comment tag
- **Episode number** from track number
- **Season number** from disc number
- **File size** for proper podcast client handling
- **MIME type** for audio format compatibility

## Server Management

### Starting the Server

```bash
# Start server on default port (8090)
chopchoprss serve

# Start server on custom port
chopchoprss serve -p 3000

# Server will display all available feeds and podcasts:
# Starting server on http://localhost:8090
# Available RSS feeds:
# - http://localhost:8090/tech-news
# Available podcast feeds:
# - http://localhost:8090/my-podcast
```

### Accessing Content

**RSS Feeds:**
```
http://localhost:8090/[feedname]
```

**Podcast Feeds:**
```
http://localhost:8090/[podcastname]        # RSS feed
http://localhost:8090/[podcastname]/audio/ # Audio files
```

## Use Cases and Workflows

### 1. Testing and Development

**Quick RSS Feed Testing:**
```bash
# Create test feed and content
chopchoprss create-feed -n test -t "Test Feed"
chopchoprss create-entry -f test -t "Test Entry" -c "This is a test"

# Start server for testing
chopchoprss serve

# Test in another terminal
curl http://localhost:8090/test
```

**Podcast Development Testing:**
```bash
# Create test podcast with sample audio directory
chopchoprss create-podcast \
  -n test-podcast \
  -t "Test Podcast" \
  -d "Testing podcast functionality" \
  -u "http://localhost:8090/test-podcast" \
  -r "./test-audio"

# Verify podcast feed
curl http://localhost:8090/test-podcast

# Test audio file serving
curl -I http://localhost:8090/test-podcast/audio/episode1.mp3
```

### 2. Production Deployment with Docker

**Single Podcast Archive:**
```bash
# Create docker-compose.yml
cat > docker-compose.yml << EOF
version: "3"
services:
  chopchoprss:
    build: .
    ports:
      - "8090:8090"
    volumes:
      - chopchoprss-data:/data
      - ./my-podcast-archive:/audio/my-podcast:ro
    restart: unless-stopped

volumes:
  chopchoprss-data:
EOF

# Build and start
docker-compose up -d

# Create podcast configuration
docker-compose exec chopchoprss create-podcast \
  -n my-podcast \
  -t "My Archived Podcast" \
  -d "Historical episodes from my podcast" \
  -u "http://your-domain.com:8090/my-podcast" \
  -r "/audio/my-podcast"

# Podcast available at http://your-domain.com:8090/my-podcast
```

**Multiple Podcast Archives (Automatic Setup):**
```bash
# Use the provided sample configuration
cp docker-compose.podcasts.yml docker-compose.override.yml

# Create startup configuration for automatic podcast setup
cp startup.json.example startup.json

# Edit startup.json to match your podcasts:
{
  "podcasts": [
    {
      "name": "show1",
      "title": "My First Show",
      "description": "Description of my show",
      "author": "Host Name",
      "email": "host@example.com",
      "baseUrl": "http://localhost:8090/show1",
      "audioDir": "/audio/show1",
      "category": "Technology"
    },
    {
      "name": "show2", 
      "title": "My Second Show",
      "description": "Another great show",
      "baseUrl": "http://localhost:8090/show2",
      "audioDir": "/audio/show2"
    }
  ]
}

# Edit docker-compose.override.yml volumes to match:
# volumes:
#   - ./podcasts/show1:/audio/show1:ro
#   - ./podcasts/show2:/audio/show2:ro
#   - ./startup.json:/data/startup.json:ro

# Start services - podcasts will be created automatically!
docker-compose up -d

# All podcasts available immediately at:
# http://localhost:8090/show1
# http://localhost:8090/show2
# http://localhost:8090 (homepage showing all feeds)
```

**Multiple Podcast Archives (Manual Setup):**
```bash
# Alternative manual approach (if you prefer docker exec commands)
cp docker-compose.podcasts.yml docker-compose.override.yml

# Edit docker-compose.override.yml volumes...
docker-compose up -d

# Configure each podcast manually
docker-compose exec chopchoprss create-podcast \
  -n show1 -t "Show 1" -d "Description" \
  -u "http://localhost:8090/show1" -r "/audio/show1"

docker-compose exec chopchoprss create-podcast \
  -n show2 -t "Show 2" -d "Description" \
  -u "http://localhost:8090/show2" -r "/audio/show2"
```

### 3. Continuous Integration / Automated Updates

**Script for Regular Updates:**
```bash
#!/bin/bash
# update-podcasts.sh - Run this when new episodes are added

# Refresh all podcasts
for podcast in $(chopchoprss list-podcasts | grep "^- " | cut -d: -f1 | cut -d' ' -f2); do
    echo "Refreshing $podcast..."
    chopchoprss refresh-podcast -n "$podcast"
done

echo "All podcasts updated!"
```

**Cron Job for Automatic Updates:**
```bash
# Add to crontab (crontab -e) to check for new episodes daily at 6 AM
0 6 * * * /path/to/update-podcasts.sh
```

### 4. Behind Reverse Proxy (Production)

**Nginx Configuration:**
```nginx
server {
    listen 80;
    server_name podcasts.example.com;

    location / {
        proxy_pass http://localhost:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Optional: Add caching for audio files
    location ~* \.(mp3|m4a|wav|flac|ogg)$ {
        proxy_pass http://localhost:8090;
        proxy_cache_valid 200 1d;
        add_header X-Cache-Status $upstream_cache_status;
    }
}
```

## Configuration

ChopChopRSS stores its configuration in JSON format:

**Default Location:** `~/.chopchoprss/config.json`

**Docker Location:** `/data/config.json` (set via `CHOPCHOP_CONFIG_DIR`)

**Startup Configuration:** `~/.chopchoprss/startup.json` or `/data/startup.json` (Docker)

**Sample Startup Configuration (startup.json):**
```json
{
  "podcasts": [
    {
      "name": "my-podcast",
      "title": "My Amazing Podcast", 
      "description": "Weekly discussions about amazing topics",
      "author": "Host Name",
      "email": "host@example.com",
      "baseUrl": "http://localhost:8090/my-podcast",
      "audioDir": "/audio/my-podcast",
      "category": "Technology",
      "language": "en",
      "imageUrl": "https://example.com/cover.jpg",
      "explicit": false
    }
  ]
}
```

**Sample Runtime Configuration Structure (config.json):**
```json
{
  "feeds": {
    "tech-news": {
      "title": "Tech News",
      "description": "Latest technology news",
      "items": [...]
    }
  },
  "podcasts": {
    "my-podcast": {
      "title": "My Podcast",
      "description": "Weekly discussions",
      "audioDir": "/audio/my-podcast",
      "baseURL": "http://localhost:8090/my-podcast",
      "episodes": [...]
    }
  }
}
```

## Troubleshooting

**Common Issues:**

1. **Audio files not found:**
   ```bash
   # Check directory permissions and paths
   ls -la /path/to/audio
   chopchoprss refresh-podcast -n podcast-name
   ```

2. **Podcast feed not updating:**
   ```bash
   # Manually refresh the podcast
   chopchoprss refresh-podcast -n podcast-name
   ```

3. **Docker volume issues:**
   ```bash
   # Check volume mounts
   docker-compose exec chopchoprss ls -la /audio/
   ```

4. **Port conflicts:**
   ```bash
   # Use different port
   chopchoprss serve -p 3000
   ```

## License

GPL3

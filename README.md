# ChopChopRSS

A fast and simple CLI tool for generating and managing RSS feeds.

## Features

- Create and manage multiple RSS feeds
- Quickly add new entries to feeds
- Support for images in feed entries
- Customizable feed metadata
- Serves feeds via HTTP
- Dockerized for easy deployment

## Installation

### From Source

```bash
git clone https://github.com/yourusername/chopchoprss.git
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

MIT

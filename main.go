// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// Config represents the application configuration
type Config struct {
	Feeds    map[string]Feed    `json:"feeds"`
	Podcasts map[string]Podcast `json:"podcasts"`
}

// Feed represents an RSS feed
type Feed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Author      string    `json:"author"`
	Email       string    `json:"email"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Items       []Item    `json:"items"`
}

// Item represents an RSS feed item
type Item struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Link        string    `json:"link"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	ImageURL    string    `json:"imageUrl,omitempty"`
}

// Podcast represents a podcast feed configuration
type Podcast struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Author      string    `json:"author"`
	Email       string    `json:"email"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	Category    string    `json:"category,omitempty"`
	Language    string    `json:"language,omitempty"`
	Copyright   string    `json:"copyright,omitempty"`
	Explicit    bool      `json:"explicit,omitempty"`
	BaseURL     string    `json:"baseUrl"`  // Base URL for serving audio files
	AudioDir    string    `json:"audioDir"` // Directory containing audio files
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Episodes    []Episode `json:"episodes"`
}

// Episode represents a podcast episode
type Episode struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	AudioURL    string        `json:"audioUrl"`
	FilePath    string        `json:"filePath"`
	Duration    time.Duration `json:"duration"`
	FileSize    int64         `json:"fileSize"`
	MimeType    string        `json:"mimeType"`
	Published   time.Time     `json:"published"`
	ImageURL    string        `json:"imageUrl,omitempty"`
	Season      int           `json:"season,omitempty"`
	Episode     int           `json:"episode,omitempty"`
}

// StartupConfig represents the startup configuration for auto-creating podcasts
type StartupConfig struct {
	Podcasts []StartupPodcast `json:"podcasts"`
}

// StartupPodcast represents a podcast configuration for auto-setup
type StartupPodcast struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link,omitempty"`
	Author      string `json:"author,omitempty"`
	Email       string `json:"email,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Category    string `json:"category,omitempty"`
	Language    string `json:"language,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Explicit    bool   `json:"explicit,omitempty"`
	BaseURL     string `json:"baseUrl"`
	AudioDir    string `json:"audioDir"`
}

var (
	cfgFile     string
	defaultPort = "8090"
	config      Config
)

func main() {
	// Create config directory if it doesn't exist
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	cfgFile = filepath.Join(configDir, "config.json")
	loadConfig()
	
	// Auto-setup podcasts from startup config if it exists
	autoSetupPodcasts(configDir)

	// Define root command
	var rootCmd = &cobra.Command{
		Use:   "chopchoprss",
		Short: "ChopChopRSS is a simple CLI tool for managing RSS feeds",
		Long:  `A CLI tool that lets you create and manage multiple RSS feeds with custom content.`,
	}

	// Create feed command
	var createFeedCmd = &cobra.Command{
		Use:   "create-feed",
		Short: "Create a new RSS feed",
		Run:   createFeed,
	}

	createFeedCmd.Flags().StringP("name", "n", "", "Feed name (required)")
	createFeedCmd.Flags().StringP("title", "t", "", "Feed title (required)")
	createFeedCmd.Flags().StringP("description", "d", "", "Feed description")
	createFeedCmd.Flags().StringP("link", "l", "", "Feed link")
	createFeedCmd.Flags().StringP("author", "a", "", "Feed author")
	createFeedCmd.Flags().StringP("email", "e", "", "Feed email")
	createFeedCmd.MarkFlagRequired("name")
	createFeedCmd.MarkFlagRequired("title")

	// Create entry command
	var createEntryCmd = &cobra.Command{
		Use:   "create-entry",
		Short: "Create a new entry in a feed",
		Run:   createEntry,
	}

	createEntryCmd.Flags().StringP("feed", "f", "", "Feed name (required)")
	createEntryCmd.Flags().StringP("title", "t", "", "Entry title (required)")
	createEntryCmd.Flags().StringP("content", "c", "", "Entry content (required)")
	createEntryCmd.Flags().StringP("link", "l", "", "Entry link")
	createEntryCmd.Flags().StringP("image", "i", "", "Entry image URL")
	createEntryCmd.MarkFlagRequired("feed")
	createEntryCmd.MarkFlagRequired("title")
	createEntryCmd.MarkFlagRequired("content")

	// List feeds command
	var listFeedsCmd = &cobra.Command{
		Use:   "list-feeds",
		Short: "List all feeds",
		Run:   listFeeds,
	}

	// Serve command
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the RSS feed server",
		Run:   serve,
	}

	serveCmd.Flags().StringP("port", "p", defaultPort, "Server port")

	// List entries command
	var listEntriesCmd = &cobra.Command{
		Use:   "list-entries",
		Short: "List all entries in a feed",
		Run:   listEntries,
	}

	listEntriesCmd.Flags().StringP("feed", "f", "", "Feed name (required)")
	listEntriesCmd.MarkFlagRequired("feed")

	// Delete feed command
	var deleteFeedCmd = &cobra.Command{
		Use:   "delete-feed",
		Short: "Delete a feed",
		Run:   deleteFeed,
	}

	deleteFeedCmd.Flags().StringP("name", "n", "", "Feed name (required)")
	deleteFeedCmd.MarkFlagRequired("name")

	// Delete entry command
	var deleteEntryCmd = &cobra.Command{
		Use:   "delete-entry",
		Short: "Delete an entry from a feed",
		Run:   deleteEntry,
	}

	deleteEntryCmd.Flags().StringP("feed", "f", "", "Feed name (required)")
	deleteEntryCmd.Flags().IntP("index", "i", -1, "Entry index (required)")
	deleteEntryCmd.MarkFlagRequired("feed")
	deleteEntryCmd.MarkFlagRequired("index")

	// Add completion command
	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(chopchoprss completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ chopchoprss completion bash > /etc/bash_completion.d/chopchoprss
  # macOS:
  $ chopchoprss completion bash > $(brew --prefix)/etc/bash_completion.d/chopchoprss

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ chopchoprss completion zsh > "${fpath[1]}/_chopchoprss"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ chopchoprss completion fish > ~/.config/fish/completions/chopchoprss.fish

PowerShell:
  PS> chopchoprss completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> chopchoprss completion powershell > chopchoprss.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	// Create podcast command
	var createPodcastCmd = &cobra.Command{
		Use:   "create-podcast",
		Short: "Create a new podcast feed from an audio directory",
		Run:   createPodcast,
	}

	createPodcastCmd.Flags().StringP("name", "n", "", "Podcast name (required)")
	createPodcastCmd.Flags().StringP("title", "t", "", "Podcast title (required)")
	createPodcastCmd.Flags().StringP("description", "d", "", "Podcast description (required)")
	createPodcastCmd.Flags().StringP("link", "l", "", "Podcast website link")
	createPodcastCmd.Flags().StringP("author", "a", "", "Podcast author/host")
	createPodcastCmd.Flags().StringP("email", "e", "", "Podcast author email")
	createPodcastCmd.Flags().StringP("image", "i", "", "Podcast cover image URL")
	createPodcastCmd.Flags().StringP("category", "c", "", "Podcast category (e.g., Technology, Comedy)")
	createPodcastCmd.Flags().StringP("language", "g", "en", "Podcast language (default: en)")
	createPodcastCmd.Flags().String("copyright", "", "Copyright information")
	createPodcastCmd.Flags().BoolP("explicit", "x", false, "Mark podcast as explicit content")
	createPodcastCmd.Flags().StringP("base-url", "u", "", "Base URL for serving audio files (required)")
	createPodcastCmd.Flags().StringP("audio-dir", "r", "", "Directory containing audio files (required)")
	createPodcastCmd.MarkFlagRequired("name")
	createPodcastCmd.MarkFlagRequired("title")
	createPodcastCmd.MarkFlagRequired("description")
	createPodcastCmd.MarkFlagRequired("base-url")
	createPodcastCmd.MarkFlagRequired("audio-dir")

	// Refresh podcast command
	var refreshPodcastCmd = &cobra.Command{
		Use:   "refresh-podcast",
		Short: "Refresh a podcast by rescanning its audio directory",
		Run:   refreshPodcast,
	}

	refreshPodcastCmd.Flags().StringP("name", "n", "", "Podcast name (required)")
	refreshPodcastCmd.MarkFlagRequired("name")

	// List podcasts command
	var listPodcastsCmd = &cobra.Command{
		Use:   "list-podcasts",
		Short: "List all podcasts",
		Run:   listPodcasts,
	}

	// Delete podcast command
	var deletePodcastCmd = &cobra.Command{
		Use:   "delete-podcast",
		Short: "Delete a podcast",
		Run:   deletePodcast,
	}

	deletePodcastCmd.Flags().StringP("name", "n", "", "Podcast name (required)")
	deletePodcastCmd.MarkFlagRequired("name")

	// Add commands to root
	rootCmd.AddCommand(createFeedCmd)
	rootCmd.AddCommand(createEntryCmd)
	rootCmd.AddCommand(listFeedsCmd)
	rootCmd.AddCommand(listEntriesCmd)
	rootCmd.AddCommand(deleteFeedCmd)
	rootCmd.AddCommand(deleteEntryCmd)
	rootCmd.AddCommand(createPodcastCmd)
	rootCmd.AddCommand(refreshPodcastCmd)
	rootCmd.AddCommand(listPodcastsCmd)
	rootCmd.AddCommand(deletePodcastCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(completionCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfigDir() string {
	// Check for environment variable first
	if envDir := os.Getenv("CHOPCHOP_CONFIG_DIR"); envDir != "" {
		return envDir
	}

	// Otherwise use home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	return filepath.Join(homeDir, ".chopchoprss")
}

func loadConfig() {
	config = Config{
		Feeds:    make(map[string]Feed),
		Podcasts: make(map[string]Podcast),
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		saveConfig()
		return
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}
}

func saveConfig() {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}

// autoSetupPodcasts checks for and processes startup configuration for auto-creating podcasts
func autoSetupPodcasts(configDir string) {
	startupConfigPath := filepath.Join(configDir, "startup.json")
	
	// Check if startup config exists
	if _, err := os.Stat(startupConfigPath); os.IsNotExist(err) {
		return // No startup config, skip
	}

	// Read startup configuration
	data, err := os.ReadFile(startupConfigPath)
	if err != nil {
		log.Printf("Warning: Failed to read startup config: %v", err)
		return
	}

	var startupConfig StartupConfig
	if err := json.Unmarshal(data, &startupConfig); err != nil {
		log.Printf("Warning: Failed to parse startup config: %v", err)
		return
	}

	// Process each podcast in startup config
	for _, podcastConfig := range startupConfig.Podcasts {
		// Skip if podcast already exists
		if _, exists := config.Podcasts[podcastConfig.Name]; exists {
			log.Printf("Podcast '%s' already exists, skipping auto-setup", podcastConfig.Name)
			continue
		}

		// Verify audio directory exists
		if _, err := os.Stat(podcastConfig.AudioDir); os.IsNotExist(err) {
			log.Printf("Warning: Audio directory '%s' for podcast '%s' does not exist, skipping", 
				podcastConfig.AudioDir, podcastConfig.Name)
			continue
		}

		// Scan audio files
		log.Printf("Auto-setting up podcast '%s' from %s...", podcastConfig.Name, podcastConfig.AudioDir)
		episodes, err := scanAudioFiles(podcastConfig.AudioDir, podcastConfig.BaseURL)
		if err != nil {
			log.Printf("Warning: Failed to scan audio files for podcast '%s': %v", 
				podcastConfig.Name, err)
			continue
		}

		// Create podcast
		now := time.Now()
		config.Podcasts[podcastConfig.Name] = Podcast{
			Title:       podcastConfig.Title,
			Description: podcastConfig.Description,
			Link:        podcastConfig.Link,
			Author:      podcastConfig.Author,
			Email:       podcastConfig.Email,
			ImageURL:    podcastConfig.ImageURL,
			Category:    podcastConfig.Category,
			Language:    podcastConfig.Language,
			Copyright:   podcastConfig.Copyright,
			Explicit:    podcastConfig.Explicit,
			BaseURL:     podcastConfig.BaseURL,
			AudioDir:    podcastConfig.AudioDir,
			Created:     now,
			Updated:     now,
			Episodes:    episodes,
		}

		log.Printf("Auto-created podcast '%s' with %d episodes", podcastConfig.Name, len(episodes))
	}

	// Save updated configuration
	if len(startupConfig.Podcasts) > 0 {
		saveConfig()
		log.Printf("Completed auto-setup of podcasts from startup configuration")
	}
}

func createFeed(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	link, _ := cmd.Flags().GetString("link")
	author, _ := cmd.Flags().GetString("author")
	email, _ := cmd.Flags().GetString("email")

	if _, exists := config.Feeds[name]; exists {
		fmt.Printf("Feed '%s' already exists\n", name)
		return
	}

	now := time.Now()
	config.Feeds[name] = Feed{
		Title:       title,
		Description: description,
		Link:        link,
		Author:      author,
		Email:       email,
		Created:     now,
		Updated:     now,
		Items:       []Item{},
	}

	saveConfig()
	fmt.Printf("Feed '%s' created successfully\n", name)
}

func createEntry(cmd *cobra.Command, args []string) {
	feedName, _ := cmd.Flags().GetString("feed")
	title, _ := cmd.Flags().GetString("title")
	content, _ := cmd.Flags().GetString("content")
	link, _ := cmd.Flags().GetString("link")
	image, _ := cmd.Flags().GetString("image")

	feed, exists := config.Feeds[feedName]
	if !exists {
		fmt.Printf("Feed '%s' does not exist\n", feedName)
		return
	}

	now := time.Now()
	newItem := Item{
		Title:       title,
		Description: content,
		Content:     content,
		Link:        link,
		Created:     now,
		Updated:     now,
		ImageURL:    image,
	}

	feed.Items = append(feed.Items, newItem)
	feed.Updated = now
	config.Feeds[feedName] = feed

	saveConfig()
	fmt.Printf("Entry '%s' added to feed '%s'\n", title, feedName)
}

func listFeeds(cmd *cobra.Command, args []string) {
	if len(config.Feeds) == 0 {
		fmt.Println("No feeds found")
		return
	}

	fmt.Println("Available feeds:")
	for name, feed := range config.Feeds {
		itemCount := len(feed.Items)
		fmt.Printf("- %s: %s (%d items)\n", name, feed.Title, itemCount)
	}
}

func listEntries(cmd *cobra.Command, args []string) {
	feedName, _ := cmd.Flags().GetString("feed")

	feed, exists := config.Feeds[feedName]
	if !exists {
		fmt.Printf("Feed '%s' does not exist\n", feedName)
		return
	}

	if len(feed.Items) == 0 {
		fmt.Printf("No entries in feed '%s'\n", feedName)
		return
	}

	fmt.Printf("Entries in feed '%s':\n", feedName)
	for i, item := range feed.Items {
		created := item.Created.Format("2006-01-02 15:04:05")
		hasImage := "no"
		if item.ImageURL != "" {
			hasImage = "yes"
		}
		fmt.Printf("[%d] %s (Created: %s, Has image: %s)\n", i, item.Title, created, hasImage)
	}
}

func deleteFeed(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	if _, exists := config.Feeds[name]; !exists {
		fmt.Printf("Feed '%s' does not exist\n", name)
		return
	}

	delete(config.Feeds, name)
	saveConfig()
	fmt.Printf("Feed '%s' deleted successfully\n", name)
}

func deleteEntry(cmd *cobra.Command, args []string) {
	feedName, _ := cmd.Flags().GetString("feed")
	index, _ := cmd.Flags().GetInt("index")

	feed, exists := config.Feeds[feedName]
	if !exists {
		fmt.Printf("Feed '%s' does not exist\n", feedName)
		return
	}

	if index < 0 || index >= len(feed.Items) {
		fmt.Printf("Invalid entry index: %d. Valid range: 0-%d\n", index, len(feed.Items)-1)
		return
	}

	// Remove the entry at the specified index
	feed.Items = append(feed.Items[:index], feed.Items[index+1:]...)
	feed.Updated = time.Now()
	config.Feeds[feedName] = feed

	saveConfig()
	fmt.Printf("Entry at index %d deleted from feed '%s'\n", index, feedName)
}

// Audio file extensions supported
var supportedAudioExts = map[string]string{
	".mp3":  "audio/mpeg",
	".m4a":  "audio/mp4",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".ogg":  "audio/ogg",
}

// scanAudioFiles scans a directory for audio files and extracts metadata
func scanAudioFiles(audioDir, baseURL string) ([]Episode, error) {
	var episodes []Episode

	err := filepath.Walk(audioDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		mimeType, supported := supportedAudioExts[ext]
		if !supported {
			return nil
		}

		// Extract metadata from audio file
		file, err := os.Open(path)
		if err != nil {
			log.Printf("Failed to open audio file %s: %v", path, err)
			return nil
		}
		defer file.Close()

		m, err := tag.ReadFrom(file)
		if err != nil {
			log.Printf("Failed to read metadata from %s: %v", path, err)
			// Continue with basic info even if metadata fails
		}

		// Get file info
		relPath, _ := filepath.Rel(audioDir, path)
		audioURL := strings.TrimSuffix(baseURL, "/") + "/audio/" + strings.ReplaceAll(relPath, "\\", "/")

		// Create episode
		episode := Episode{
			Title:     getStringOrDefault(m, "title", filepath.Base(path)),
			AudioURL:  audioURL,
			FilePath:  path,
			FileSize:  info.Size(),
			MimeType:  mimeType,
			Published: info.ModTime(),
		}

		if m != nil {
			episode.Description = getStringOrDefault(m, "comment", "")
			if album := m.Album(); album != "" {
				episode.Description = album + " - " + episode.Description
			}

			// Try to extract episode/season numbers
			if track, total := m.Track(); track != 0 {
				episode.Episode = track
				_ = total // Could be used for validation
			}

			// Try to extract season from album or genre
			if disc, _ := m.Disc(); disc != 0 {
				episode.Season = disc
			}
		}

		episodes = append(episodes, episode)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan audio directory: %v", err)
	}

	// Sort episodes by published date (oldest first)
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].Published.Before(episodes[j].Published)
	})

	return episodes, nil
}

// Helper function to safely get string values from metadata
func getStringOrDefault(m tag.Metadata, field, defaultValue string) string {
	if m == nil {
		return defaultValue
	}

	switch field {
	case "title":
		if title := m.Title(); title != "" {
			return title
		}
	case "artist":
		if artist := m.Artist(); artist != "" {
			return artist
		}
	case "album":
		if album := m.Album(); album != "" {
			return album
		}
	case "comment":
		if comment := m.Comment(); comment != "" {
			return comment
		}
	}

	return defaultValue
}

// createPodcast creates a new podcast feed from a directory of audio files
func createPodcast(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	title, _ := cmd.Flags().GetString("title")
	description, _ := cmd.Flags().GetString("description")
	link, _ := cmd.Flags().GetString("link")
	author, _ := cmd.Flags().GetString("author")
	email, _ := cmd.Flags().GetString("email")
	imageURL, _ := cmd.Flags().GetString("image")
	category, _ := cmd.Flags().GetString("category")
	language, _ := cmd.Flags().GetString("language")
	copyright, _ := cmd.Flags().GetString("copyright")
	explicit, _ := cmd.Flags().GetBool("explicit")
	baseURL, _ := cmd.Flags().GetString("base-url")
	audioDir, _ := cmd.Flags().GetString("audio-dir")

	if _, exists := config.Podcasts[name]; exists {
		fmt.Printf("Podcast '%s' already exists\n", name)
		return
	}

	// Verify audio directory exists
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		fmt.Printf("Audio directory '%s' does not exist\n", audioDir)
		return
	}

	// Scan audio files
	fmt.Printf("Scanning audio files in %s...\n", audioDir)
	episodes, err := scanAudioFiles(audioDir, baseURL)
	if err != nil {
		fmt.Printf("Failed to scan audio files: %v\n", err)
		return
	}

	now := time.Now()
	config.Podcasts[name] = Podcast{
		Title:       title,
		Description: description,
		Link:        link,
		Author:      author,
		Email:       email,
		ImageURL:    imageURL,
		Category:    category,
		Language:    language,
		Copyright:   copyright,
		Explicit:    explicit,
		BaseURL:     baseURL,
		AudioDir:    audioDir,
		Created:     now,
		Updated:     now,
		Episodes:    episodes,
	}

	saveConfig()
	fmt.Printf("Podcast '%s' created with %d episodes\n", name, len(episodes))
}

// refreshPodcast rescans the audio directory and updates episodes
func refreshPodcast(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	podcast, exists := config.Podcasts[name]
	if !exists {
		fmt.Printf("Podcast '%s' does not exist\n", name)
		return
	}

	// Rescan audio files
	fmt.Printf("Rescanning audio files in %s...\n", podcast.AudioDir)
	episodes, err := scanAudioFiles(podcast.AudioDir, podcast.BaseURL)
	if err != nil {
		fmt.Printf("Failed to scan audio files: %v\n", err)
		return
	}

	podcast.Episodes = episodes
	podcast.Updated = time.Now()
	config.Podcasts[name] = podcast

	saveConfig()
	fmt.Printf("Podcast '%s' refreshed with %d episodes\n", name, len(episodes))
}

// listPodcasts lists all configured podcasts
func listPodcasts(cmd *cobra.Command, args []string) {
	if len(config.Podcasts) == 0 {
		fmt.Println("No podcasts found")
		return
	}

	fmt.Println("Available podcasts:")
	for name, podcast := range config.Podcasts {
		episodeCount := len(podcast.Episodes)
		fmt.Printf("- %s: %s (%d episodes)\n", name, podcast.Title, episodeCount)
	}
}

// deletePodcast removes a podcast
func deletePodcast(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")

	if _, exists := config.Podcasts[name]; !exists {
		fmt.Printf("Podcast '%s' does not exist\n", name)
		return
	}

	delete(config.Podcasts, name)
	saveConfig()
	fmt.Printf("Podcast '%s' deleted successfully\n", name)
}

func serve(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")

	r := mux.NewRouter()

	// Handle homepage
	r.HandleFunc("/", serveHomepage)
	
	// Serve the logo
	r.HandleFunc("/chopchop.png", serveLogo)

	// Handle regular RSS feeds
	for name := range config.Feeds {
		feedName := name // Capture for closure
		r.HandleFunc("/"+feedName, func(w http.ResponseWriter, r *http.Request) {
			serveRSSFeed(w, feedName)
		})
	}

	// Handle podcast feeds
	for name := range config.Podcasts {
		podcastName := name // Capture for closure
		r.HandleFunc("/"+podcastName, func(w http.ResponseWriter, r *http.Request) {
			servePodcastFeed(w, podcastName)
		})
	}

	// Serve audio files
	for name, podcast := range config.Podcasts {
		podcastName := name // Capture for closure
		audioDir := podcast.AudioDir
		r.PathPrefix("/" + podcastName + "/audio/").Handler(
			http.StripPrefix("/"+podcastName+"/audio/", http.FileServer(http.Dir(audioDir))),
		)
	}

	fmt.Printf("Starting server on http://localhost:%s\n", port)

	if len(config.Feeds) > 0 {
		fmt.Println("Available RSS feeds:")
		for name := range config.Feeds {
			fmt.Printf("- http://localhost:%s/%s\n", port, name)
		}
	}

	if len(config.Podcasts) > 0 {
		fmt.Println("Available podcast feeds:")
		for name := range config.Podcasts {
			fmt.Printf("- http://localhost:%s/%s\n", port, name)
		}
	}

	if len(config.Feeds) == 0 && len(config.Podcasts) == 0 {
		fmt.Println("No feeds or podcasts configured")
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func serveRSSFeed(w http.ResponseWriter, feedName string) {
	feed, exists := config.Feeds[feedName]
	if !exists {
		http.NotFound(w, nil)
		return
	}

	// Convert our feed structure to gorilla/feeds format
	// now := time.Now()
	f := &feeds.Feed{
		Title:       feed.Title,
		Link:        &feeds.Link{Href: feed.Link},
		Description: feed.Description,
		Author:      &feeds.Author{Name: feed.Author, Email: feed.Email},
		Created:     feed.Created,
		Updated:     feed.Updated,
	}

	f.Items = make([]*feeds.Item, len(feed.Items))
	for i, item := range feed.Items {
		feedItem := &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
			Content:     item.Content,
			Created:     item.Created,
			Updated:     item.Updated,
		}

		if item.ImageURL != "" {
			feedItem.Enclosure = &feeds.Enclosure{
				Url:    item.ImageURL,
				Length: "0",
				Type:   "image/jpeg", // Assuming JPEG, but could be determined based on extension
			}
		}

		f.Items[i] = feedItem
	}

	rss, err := f.ToRss()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(rss))
}

// servePodcastFeed serves a podcast feed as RSS with podcast-specific elements
func servePodcastFeed(w http.ResponseWriter, podcastName string) {
	podcast, exists := config.Podcasts[podcastName]
	if !exists {
		http.NotFound(w, nil)
		return
	}

	// Convert our podcast structure to gorilla/feeds format
	f := &feeds.Feed{
		Title:       podcast.Title,
		Link:        &feeds.Link{Href: podcast.Link},
		Description: podcast.Description,
		Author:      &feeds.Author{Name: podcast.Author, Email: podcast.Email},
		Created:     podcast.Created,
		Updated:     podcast.Updated,
		Copyright:   podcast.Copyright,
	}

	// Add podcast-specific image
	if podcast.ImageURL != "" {
		f.Image = &feeds.Image{
			Url:   podcast.ImageURL,
			Title: podcast.Title,
			Link:  podcast.Link,
		}
	}

	// Convert episodes to feed items
	f.Items = make([]*feeds.Item, len(podcast.Episodes))
	for i, episode := range podcast.Episodes {
		feedItem := &feeds.Item{
			Title:       episode.Title,
			Description: episode.Description,
			Created:     episode.Published,
			Updated:     episode.Published,
		}

		// Add audio enclosure for podcast episode
		feedItem.Enclosure = &feeds.Enclosure{
			Url:    episode.AudioURL,
			Length: strconv.FormatInt(episode.FileSize, 10),
			Type:   episode.MimeType,
		}

		// Add episode image if available
		if episode.ImageURL != "" {
			feedItem.Description += fmt.Sprintf(`<br><img src="%s" alt="Episode Image">`, episode.ImageURL)
		}

		// Add episode and season numbers to description if available
		if episode.Season > 0 || episode.Episode > 0 {
			episodeInfo := ""
			if episode.Season > 0 {
				episodeInfo += fmt.Sprintf("Season %d", episode.Season)
			}
			if episode.Episode > 0 {
				if episodeInfo != "" {
					episodeInfo += ", "
				}
				episodeInfo += fmt.Sprintf("Episode %d", episode.Episode)
			}
			feedItem.Description = fmt.Sprintf("[%s] %s", episodeInfo, feedItem.Description)
		}

		f.Items[i] = feedItem
	}

	// Generate RSS with podcast extensions
	rss, err := f.ToRss()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add podcast-specific XML namespaces and elements
	rss = addPodcastExtensions(rss, podcast)

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(rss))
}

// addPodcastExtensions adds iTunes and other podcast-specific XML elements
func addPodcastExtensions(rss string, podcast Podcast) string {
	// Add iTunes namespace
	rss = strings.Replace(rss, "<rss version=\"2.0\">",
		`<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">`, 1)

	// Find the channel opening tag and add iTunes elements after it
	channelStart := strings.Index(rss, "<channel>")
	if channelStart == -1 {
		return rss
	}

	insertPos := channelStart + len("<channel>")

	var itunesElements strings.Builder

	// Add iTunes category
	if podcast.Category != "" {
		itunesElements.WriteString(fmt.Sprintf("\n    <itunes:category text=\"%s\" />", podcast.Category))
	}

	// Add iTunes explicit tag
	explicitValue := "no"
	if podcast.Explicit {
		explicitValue = "yes"
	}
	itunesElements.WriteString(fmt.Sprintf("\n    <itunes:explicit>%s</itunes:explicit>", explicitValue))

	// Add iTunes author
	if podcast.Author != "" {
		itunesElements.WriteString(fmt.Sprintf("\n    <itunes:author>%s</itunes:author>", podcast.Author))
	}

	// Add iTunes owner
	if podcast.Email != "" || podcast.Author != "" {
		itunesElements.WriteString("\n    <itunes:owner>")
		if podcast.Author != "" {
			itunesElements.WriteString(fmt.Sprintf("\n      <itunes:name>%s</itunes:name>", podcast.Author))
		}
		if podcast.Email != "" {
			itunesElements.WriteString(fmt.Sprintf("\n      <itunes:email>%s</itunes:email>", podcast.Email))
		}
		itunesElements.WriteString("\n    </itunes:owner>")
	}

	// Add iTunes image
	if podcast.ImageURL != "" {
		itunesElements.WriteString(fmt.Sprintf("\n    <itunes:image href=\"%s\" />", podcast.ImageURL))
	}

	// Add language
	if podcast.Language != "" {
		itunesElements.WriteString(fmt.Sprintf("\n    <language>%s</language>", podcast.Language))
	}

	// Insert the iTunes elements
	result := rss[:insertPos] + itunesElements.String() + rss[insertPos:]

	return result
}

// serveHomepage serves a nice HTML homepage with logo and feed information
func serveHomepage(w http.ResponseWriter, r *http.Request) {
	// Only serve homepage on exact root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	feedCount := len(config.Feeds)
	podcastCount := len(config.Podcasts)
	
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ChopChopRSS</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f8f9fa;
        }
        .container {
            background: white;
            border-radius: 10px;
            padding: 40px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
        }
        .logo {
            max-width: 200px;
            height: auto;
            margin-bottom: 20px;
        }
        h1 {
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .subtitle {
            color: #7f8c8d;
            font-size: 1.2em;
            margin-bottom: 30px;
        }
        .stats {
            display: flex;
            justify-content: center;
            gap: 40px;
            margin: 30px 0;
            flex-wrap: wrap;
        }
        .stat {
            text-align: center;
            padding: 20px;
            background: #ecf0f1;
            border-radius: 8px;
            min-width: 120px;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #3498db;
            display: block;
        }
        .stat-label {
            color: #7f8c8d;
            text-transform: uppercase;
            font-size: 0.9em;
            letter-spacing: 1px;
        }
        .feeds-section {
            margin-top: 40px;
        }
        .feeds-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
            margin-top: 20px;
        }
        .feed-list {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
        }
        .feed-list h3 {
            color: #2c3e50;
            margin-top: 0;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .feed-item {
            margin: 10px 0;
            padding: 8px 0;
            border-bottom: 1px solid #ecf0f1;
        }
        .feed-item:last-child {
            border-bottom: none;
        }
        .feed-link {
            color: #3498db;
            text-decoration: none;
            font-weight: 500;
        }
        .feed-link:hover {
            text-decoration: underline;
        }
        .feed-description {
            color: #7f8c8d;
            font-size: 0.9em;
            margin-top: 5px;
        }
        .no-feeds {
            text-align: center;
            color: #7f8c8d;
            font-style: italic;
            padding: 40px 20px;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
            color: #7f8c8d;
            font-size: 0.9em;
        }
        .icon {
            width: 20px;
            height: 20px;
            display: inline-block;
        }
        @media (max-width: 600px) {
            .stats {
                gap: 20px;
            }
            .stat {
                min-width: 100px;
                padding: 15px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="/chopchop.png" alt="ChopChopRSS Logo" class="logo">
            <h1>ChopChopRSS</h1>
            <div class="subtitle">Fast and simple RSS feeds and podcast hosting</div>
        </div>

        <div class="stats">
            <div class="stat">
                <span class="stat-number">` + fmt.Sprintf("%d", feedCount) + `</span>
                <span class="stat-label">RSS Feeds</span>
            </div>
            <div class="stat">
                <span class="stat-number">` + fmt.Sprintf("%d", podcastCount) + `</span>
                <span class="stat-label">Podcasts</span>
            </div>
        </div>`

	if feedCount > 0 || podcastCount > 0 {
		html += `
        <div class="feeds-section">
            <div class="feeds-grid">`
		
		if feedCount > 0 {
			html += `
                <div class="feed-list">
                    <h3>
                        <span class="icon">📰</span>
                        RSS Feeds
                    </h3>`
			for name, feed := range config.Feeds {
				itemCount := len(feed.Items)
				html += fmt.Sprintf(`
                    <div class="feed-item">
                        <a href="/%s" class="feed-link">%s</a>
                        <div class="feed-description">%s • %d items</div>
                    </div>`, name, feed.Title, feed.Description, itemCount)
			}
			html += `
                </div>`
		}

		if podcastCount > 0 {
			html += `
                <div class="feed-list">
                    <h3>
                        <span class="icon">🎧</span>
                        Podcast Feeds
                    </h3>`
			for name, podcast := range config.Podcasts {
				episodeCount := len(podcast.Episodes)
				html += fmt.Sprintf(`
                    <div class="feed-item">
                        <a href="/%s" class="feed-link">%s</a>
                        <div class="feed-description">%s • %d episodes</div>
                    </div>`, name, podcast.Title, podcast.Description, episodeCount)
			}
			html += `
                </div>`
		}

		html += `
            </div>
        </div>`
	} else {
		html += `
        <div class="no-feeds">
            <h3>No feeds or podcasts configured yet</h3>
            <p>Use the ChopChopRSS CLI to create your first RSS feed or podcast.</p>
        </div>`
	}

	html += `
        <div class="footer">
            <p>Powered by <strong>ChopChopRSS</strong> • 
            <a href="https://github.com/madeofpendletonwool/chopchoprss" style="color: #3498db;">GitHub</a></p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

// serveLogo serves the chopchop.png logo file
func serveLogo(w http.ResponseWriter, r *http.Request) {
	// Try to serve the logo from the current directory
	logoPath := "chopchop.png"
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		// If not found, return a simple 404
		http.NotFound(w, r)
		return
	}
	
	http.ServeFile(w, r, logoPath)
}

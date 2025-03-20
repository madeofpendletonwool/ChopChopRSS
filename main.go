// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// Config represents the application configuration
type Config struct {
	Feeds map[string]Feed `json:"feeds"`
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

	// Add commands to root
	rootCmd.AddCommand(createFeedCmd)
	rootCmd.AddCommand(createEntryCmd)
	rootCmd.AddCommand(listFeedsCmd)
	rootCmd.AddCommand(listEntriesCmd)
	rootCmd.AddCommand(deleteFeedCmd)
	rootCmd.AddCommand(deleteEntryCmd)
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
		Feeds: make(map[string]Feed),
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		saveConfig()
		return
	}

	data, err := ioutil.ReadFile(cfgFile)
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

	if err := ioutil.WriteFile(cfgFile, data, 0644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
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

func serve(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetString("port")

	r := mux.NewRouter()

	for name := range config.Feeds {
		feedName := name // Capture for closure
		r.HandleFunc("/"+feedName, func(w http.ResponseWriter, r *http.Request) {
			serveRSSFeed(w, feedName)
		})
	}

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	fmt.Println("Available feeds:")
	for name := range config.Feeds {
		fmt.Printf("- http://localhost:%s/%s\n", port, name)
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

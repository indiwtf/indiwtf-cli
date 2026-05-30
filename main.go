package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// version is the current release of the CLI. It is overridden at build
// time via -ldflags "-X main.version=<tag>".
var version = "dev"

// repoSlug is the GitHub owner/repo used for updates.
const repoSlug = "indiwtf/indiwtf-cli"

// domainRegex matches a valid fully-qualified domain name.
var domainRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// isValidDomain reports whether the given hostname is a syntactically valid domain.
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	return domainRegex.MatchString(domain)
}

type stringSliceFlag []string

func (f *stringSliceFlag) String() string {
	return strings.Join(*f, ", ")
}

func (f *stringSliceFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type DomainStatus struct {
	Domain string `json:"domain"`
	Status string `json:"status"`
	IP     string `json:"ip"`
}

type Config struct {
	Token string `json:"token"`
}

var token string
var configFilePath string

func init() {
	// Define the path to the configuration file in the user's home directory
	configFilePath = getHomeDir() + "/.indiwtf/config.json"

	// Load the API token from the configuration file, if available.
	config := loadConfig()
	token = config.Token
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		os.Exit(1)
	}
	return home
}

// loadConfig loads the API token from a configuration file.
func loadConfig() Config {
	config := Config{}
	file, err := os.Open(configFilePath)
	if err != nil {
		// If the file doesn't exist, return an empty configuration.
		return config
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		// If there is an error decoding the file, return an empty configuration.
		return config
	}

	return config
}

// saveConfig saves the API token to a configuration file.
func saveConfig(config Config) error {
	// Ensure the directory exists
	configDir := getHomeDir() + "/.indiwtf"
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

// auth stores the given API token in the configuration file.
func auth(apiToken string) {
	apiToken = strings.TrimSpace(apiToken)
	if apiToken == "" {
		fmt.Println("Usage: indiwtf auth API_TOKEN")
		fmt.Println("Get an API token at https://indiwtf.com/pricing")
		return
	}

	if err := saveConfig(Config{Token: apiToken}); err != nil {
		fmt.Printf("Error saving the API token to the configuration file: %v\n", err)
		return
	}

	fmt.Printf("API token saved to %s\n", configFilePath)
}

// uninstall removes the indiwtf binary and the configuration directory.
func uninstall() {
	// Remove the configuration directory (~/.indiwtf).
	configDir := getHomeDir() + "/.indiwtf"
	if err := os.RemoveAll(configDir); err != nil {
		fmt.Printf("Error removing configuration directory %s: %v\n", configDir, err)
	} else {
		fmt.Printf("Removed configuration directory: %s\n", configDir)
	}

	// Resolve the path to the running binary and remove it.
	binPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error locating the indiwtf binary: %v\n", err)
		return
	}

	if err := os.Remove(binPath); err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Permission denied removing %s.\n", binPath)
			fmt.Printf("Try again with elevated privileges: sudo rm %s\n", binPath)
		} else {
			fmt.Printf("Error removing the indiwtf binary %s: %v\n", binPath, err)
		}
		return
	}

	fmt.Printf("Removed binary: %s\n", binPath)
	fmt.Println("Indiwtf CLI has been uninstalled.")
}

// githubRelease holds the fields we need from the GitHub releases API.
type githubRelease struct {
	TagName string `json:"tag_name"`
}

// latestVersion queries the GitHub API for the latest released tag.
func latestVersion() (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repoSlug)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "indiwtf-cli/"+version)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response from GitHub: %s", resp.Status)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("no released version found")
	}
	return release.TagName, nil
}

// update downloads the latest released binary and replaces the running one.
func update() {
	fmt.Println("Checking for the latest version...")
	tag, err := latestVersion()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		return
	}

	if tag == version {
		fmt.Printf("You are already on the latest version (%s).\n", version)
		return
	}

	fmt.Printf("Updating from %s to %s...\n", version, tag)

	binPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error locating the indiwtf binary: %v\n", err)
		return
	}
	// Resolve symlinks so we replace the real binary, not a link to it.
	if resolved, err := filepath.EvalSymlinks(binPath); err == nil {
		binPath = resolved
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/indiwtf", repoSlug, tag)
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		fmt.Printf("Error creating download request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "indiwtf-cli/"+version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error downloading the update: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error downloading the update: unexpected response %s\n", resp.Status)
		return
	}

	// Write to a temporary file in the same directory for an atomic replace.
	tmpFile, err := os.CreateTemp(filepath.Dir(binPath), ".indiwtf-update-*")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v\n", err)
		if os.IsPermission(err) {
			fmt.Printf("Try again with elevated privileges: sudo %s update\n", binPath)
		}
		return
	}
	tmpPath := tmpFile.Name()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		fmt.Printf("Error saving the update: %v\n", err)
		return
	}
	tmpFile.Close()

	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		fmt.Printf("Error setting permissions: %v\n", err)
		return
	}

	if err := os.Rename(tmpPath, binPath); err != nil {
		os.Remove(tmpPath)
		fmt.Printf("Error replacing the binary %s: %v\n", binPath, err)
		if os.IsPermission(err) {
			fmt.Printf("Try again with elevated privileges: sudo %s update\n", binPath)
		}
		return
	}

	fmt.Printf("Updated to %s.\n", tag)
}

// checkDomain sends an HTTP GET request to the API endpoint with the token and returns the status and IP of the domain.
func checkDomain(domain string) (*DomainStatus, error) {
	if token == "" {
		fmt.Println("API token is required. Please enter your API token (https://indiwtf.com/pricing):")
		fmt.Scanln(&token)
		config := Config{
			Token: token,
		}
		err := saveConfig(config)
		if err != nil {
			fmt.Printf("Error saving the API token to the configuration file: %v\n", err)
		}
	}

	apiURL := fmt.Sprintf("https://indiwtf.com/api/check?domain=%s&token=%s", url.QueryEscape(domain), token)

	// Create an HTTP client with a custom User-Agent string
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Set a custom User-Agent string
	req.Header.Set("User-Agent", "indiwtf-cli/"+version)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var domainStatus DomainStatus
	err = json.NewDecoder(resp.Body).Decode(&domainStatus)
	if err != nil {
		return nil, err
	}

	return &domainStatus, nil
}

func main() {
	// Instructions for running the program.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [domain1] [domain2] ...\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Check if one or more domains are blocked in Indonesia.")
		fmt.Fprintln(os.Stderr, "\nCommands:")
		fmt.Fprintln(os.Stderr, "  auth TOKEN   Save your API token to the configuration file")
		fmt.Fprintln(os.Stderr, "  update       Update indiwtf to the latest version")
		fmt.Fprintln(os.Stderr, "  uninstall    Remove the indiwtf binary and configuration files")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		fmt.Fprintln(os.Stderr, "  -h, --help   Show this help message and exit")
		flag.PrintDefaults()
	}

	// Parse command-line flags.
	flag.Parse()

	args := flag.Args()

	// If no arguments are provided, show the usage and exit.
	if len(args) == 0 {
		flag.Usage()
		return
	}

	// Handle subcommands.
	switch args[0] {
	case "auth":
		auth(strings.Join(args[1:], " "))
		return
	case "update":
		update()
		return
	case "uninstall":
		uninstall()
		return
	}

	// Iterate over the domain names and perform the necessary checks.
	for _, rawURL := range args {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			continue
		}

		// If the scheme is empty, assume HTTPS as the default scheme.
		if parsedURL.Scheme == "" {
			rawURL = "https://" + rawURL
			parsedURL, err = url.Parse(rawURL)
			if err != nil {
				fmt.Printf("Error parsing URL: %v\n", err)
				continue
			}
		}

		hostname := parsedURL.Hostname()

		// Skip inputs that are not valid domain names.
		if !isValidDomain(hostname) {
			fmt.Printf("Invalid domain: %s\n", rawURL)
			continue
		}

		domainStatus, err := checkDomain(hostname)
		if err != nil {
			fmt.Printf("Error checking domain %s: %v\n", hostname, err)
			continue
		}

		fmt.Printf("Domain: %s | Status: %s | IP: %s\n", domainStatus.Domain, domainStatus.Status, domainStatus.IP)
	}
}

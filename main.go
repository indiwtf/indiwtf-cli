package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

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
	req.Header.Set("User-Agent", "indiwtf-cli/1.0")

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
		fmt.Fprintf(os.Stderr, "Usage: %s [domain1] [domain2] ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse command-line flags.
	flag.Parse()

	domains := flag.Args()

	// If no domain names are provided, show the usage and exit.
	if len(domains) == 0 {
		flag.Usage()
		return
	}

	// Iterate over the domain names and perform the necessary checks.
	for _, rawURL := range domains {
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

		domainStatus, err := checkDomain(parsedURL.Hostname())
		if err != nil {
			fmt.Printf("Error checking domain %s: %v\n", parsedURL.Hostname(), err)
			continue
		}

		fmt.Printf("Domain: %s | Status: %s | IP: %s\n", domainStatus.Domain, domainStatus.Status, domainStatus.IP)
	}
}

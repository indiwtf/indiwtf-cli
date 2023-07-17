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

// checkDomain sends an HTTP GET request to the API endpoint and returns the status and IP of the domain.
func checkDomain(domain string) (*DomainStatus, error) {
	apiURL := fmt.Sprintf("https://indiwtf.upset.dev/api/check?domain=%s", url.QueryEscape(domain))

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

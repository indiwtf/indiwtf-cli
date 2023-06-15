package main

import (
	"context"
	"flag"
	"fmt"
	"net"
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

// resolveDNS resolves the IP address for a given hostname using the specified DNS server.
func resolveDNS(hostname, dnsServer string) (string, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{}
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}

	// Use LookupIPAddr to perform DNS resolution and get IP addresses associated with the hostname.
	ips, err := resolver.LookupIPAddr(context.Background(), hostname)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for %s", hostname)
	}

	// Return the first IP address as a string representation.
	return ips[0].IP.String(), nil
}

func main() {
	var domainFlags stringSliceFlag

	// Customize usage message to provide instructions for running the program.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [domain1] [domain2] ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Define a custom flag for specifying domain names.
	flag.Var(&domainFlags, "d", "Domain names")

	// Parse command-line flags.
	flag.Parse()

	domains := flag.Args()

	// If domain names are provided using the -d flag, add them to the list of domains.
	if len(domainFlags) > 0 {
		domains = append(domains, domainFlags...)
	}

	// If no domain names are provided, show the usage and exit.
	if len(domains) == 0 {
		flag.Usage()
		return
	}

	// Currently we are using the Telkom DNS directly, if we want to cache queries locally we will need to proxy it in the future.
	dnsServer := "118.98.44.10:53"

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

		ipAddress, err := resolveDNS(parsedURL.Hostname(), dnsServer)
		if err != nil {
			fmt.Printf("Error resolving IP address for %s: %v\n", rawURL, err)
			continue
		}

		status := "Not Blocked"
		if ipAddress == "36.86.63.185" {
			status = "Blocked"
		}

		fmt.Printf("Domain: %s | Status: %s\n", parsedURL.Hostname(), status)
	}
}

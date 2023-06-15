# Indiwtf CLI - Blocked Website Checker for Indonesia

Indiwtf CLI is a tiny command-line tool written in Go that allows you to check if your website is blocked in Indonesia. It performs DNS resolution using a specified DNS server in Indonesia and checks the accessibility status of the website based on the resolved IP address.

For web version you can visit at [indiwtf.upset.dev](https://indiwtf.upset.dev).

## Features

- Check the accessibility status of a website based on the resolved IP address.
- Resolve IP address for a given hostname using a custom DNS server in Indonesia.
- Supports checking multiple websites in a single run.

## Prerequisites

Go (version 1.16 or above).

## Installation

1. Clone the repository:

```
git clone https://github.com/fransallen/indiwtf-cli.git
```

2. Change to the project directory:

```
cd indiwtf-cli
```

3. Build the executable:

```
./make
```

## Usage

Run the `indiwtf` executable with the desired options:

```
./indiwtf [flags] [domain1] [domain2] ...
```

## Flags

- **-d**: Specify a domain name to check. You can provide multiple domain names by repeating the flag.

## Examples

Check the accessibility status of a single website:

```
./indiwtf example.com
```

Check the accessibility status of multiple websites:

```
./indiwtf -d puredns.org -d github.com -d reddit.com
```

## Customization

You can customize the default DNS server by modifying the `dnsServer` variable in the code. By default, it uses the Telkom DNS.

## Limitations

Please note that Indiwtf specifically focuses on checking website accessibility in Indonesia. The results may not accurately represent the accessibility status of websites in other regions.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project was inspired by the need to quickly check if websites are blocked in Indonesia.

## Contributing

Contributions are welcome! If you find any issues or want to add new features, please open an issue or submit a pull request.

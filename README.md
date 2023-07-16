# Indiwtf CLI

Will your website be blocked by Kominfo? Let's check!

Indiwtf CLI is a tiny command-line tool written in Go that allows you to check if your website is blocked in Indonesia. It performs DNS resolution using a specified DNS server in Indonesia and checks the accessibility status of the website based on the resolved IP address.

You can access the web version by visiting [indiwtf.upset.dev](https://indiwtf.upset.dev). Indiwtf is also available in a [Telegram Bot](https://github.com/fransallen/indiwtf-telegram-bot) version.

## Usage

Run the `indiwtf` executable with the desired options:

```
./indiwtf [domain1] [domain2] ...
```

## Examples

Check the accessibility status of a single website:

```
./indiwtf example.com
```

Check the accessibility status of multiple websites:

```
./indiwtf puredns.org github.com reddit.com
```

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

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project was inspired by the need to quickly check if websites are blocked in Indonesia.

## Contributing

Contributions are welcome! If you find any issues or want to add new features, please open an issue or submit a pull request.

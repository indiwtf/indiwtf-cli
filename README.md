# Indiwtf CLI

Will your website be blocked by Kominfo? Let's check!

Indiwtf CLI is a tiny command-line tool written in Go that allows you to check if your website is blocked in Indonesia. It uses the [Indiwtf API](https://indiwtf.com/api) to check the accessibility status of the website in Indonesian Internet network.

You can access the web version by visiting [indiwtf.com](https://indiwtf.com). Indiwtf is also available in a [Telegram Bot](https://github.com/fransallen/indiwtf-telegram-bot) version.

## Usage

Run `indiwtf` with one or more domains:

```sh
indiwtf [domain1] [domain2] ...
```

## Examples

Check the accessibility status of a single website:

```sh
indiwtf example.com
```

Check the accessibility status of multiple websites:

```sh
indiwtf puredns.org github.com reddit.com
```

## API Token

Indiwtf CLI requires an API token to check website accessibility. You can store the API token securely in a configuration file located at `~/.indiwtf/config.json`.

Save your token with the `auth` command:

```sh
indiwtf auth API_TOKEN
```

Alternatively, the program will prompt you to enter the token the first time you run a check if it's not found in the configuration file.

To obtain an API token, please visit [indiwtf.com/pricing](https://indiwtf.com/pricing).

## Installation

Install the latest prebuilt binary with a single command:

```sh
curl -fsSL https://github.com/indiwtf/indiwtf-cli/raw/main/install.sh | sh
```

Or using `wget`:

```sh
wget -qO- https://github.com/indiwtf/indiwtf-cli/raw/main/install.sh | sh
```

The script downloads the `indiwtf` binary to `/usr/local/bin` and makes it executable. To install somewhere else, set `INSTALL_DIR`:

```sh
INSTALL_DIR="$HOME/.local/bin" sh -c "$(curl -fsSL https://github.com/indiwtf/indiwtf-cli/raw/main/install.sh)"
```

<details>
<summary>Manual installation</summary>

Download the latest binary from the [releases page](https://github.com/indiwtf/indiwtf-cli/releases/latest) and place it under `/usr/local/bin`:

```sh
sudo wget -O /usr/local/bin/indiwtf \
  https://github.com/indiwtf/indiwtf-cli/releases/latest/download/indiwtf
sudo chmod +x /usr/local/bin/indiwtf
```

For other platforms, replace the asset name (e.g. `indiwtf-darwin-arm64`, `indiwtf-linux-arm64`, `indiwtf-windows-amd64.exe`).

</details>

### Windows

Prebuilt Windows binaries are published on the [releases page](https://github.com/indiwtf/indiwtf-cli/releases/latest) (`indiwtf-windows-amd64.exe` and `indiwtf-windows-arm64.exe`). Download one, rename it to `indiwtf.exe`, and place it in a directory on your `PATH`.

> The `install.sh` script, `indiwtf update`, and `indiwtf uninstall` are not available on Windows. Update by downloading the latest `.exe`, and uninstall by deleting the file.

## Update

Update to the latest release at any time:

```sh
indiwtf update
```

## Uninstall

Remove the **indiwtf** binary and its configuration files (`~/.indiwtf`):

```sh
indiwtf uninstall
```

<details>
<summary>Manual uninstall</summary>

```sh
sudo rm /usr/local/bin/indiwtf
rm -rf ~/.indiwtf
```

</details>

## Features

- Check the accessibility status of a website based on the resolved IP address.
- Resolve IP address for a given hostname using a custom DNS server in Indonesia.
- Supports checking multiple websites in a single run.

## Development

You will need Go (version 1.16 or above).

Clone the repository:

```sh
git clone https://github.com/fransallen/indiwtf-cli.git
cd indiwtf-cli
```

To build the binary, run:

```sh
CGO_ENABLED=0 go build -o indiwtf main.go
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! If you find any issues or want to add new features, please open an issue or submit a pull request.

# mDNS Browser

A modern, terminal-based mDNS/Bonjour service browser built with Go. Discover and explore network services in your local network with an elegant, keyboard-driven interface.

## Features

- **Real-time Service Discovery**: Automatically detects mDNS services as they appear on your network
- **570+ Service Types**: Supports a comprehensive list of mDNS service types including HTTP, SSH, AirPlay, printers, and many more
- **Split-Pane Interface**: Browse services in the left pane while viewing detailed information in the right pane
- **Rich Service Details**: View service names, hostnames, IPv4/IPv6 addresses, ports, and additional metadata
- **Keyboard Navigation**: Vim-style keybindings for efficient navigation
- **Context-Aware Help**: Dynamic help system that shows relevant commands based on your current focus
- **Graceful Shutdown**: Clean exit with proper signal handling

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/yourusername/mdns-browser/releases).

Available platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/mdns-browser.git
cd mdns-browser

# Build and install
make build

# Or install directly with Go
go install ./cmd/mdns-browser
```

### Using Go Install

```bash
go install github.com/yourusername/mdns-browser/cmd/mdns-browser@latest
```

## Usage

Simply run the binary to start discovering services:

```bash
mdns-browser
```

### Keyboard Shortcuts

#### Common
- `q` or `Ctrl+C` - Quit the application
- `Tab` - Switch focus between service list and details pane
- `?` - Toggle help view (short/full)

#### Service List (left pane)
- `↑`/`k` - Move up
- `↓`/`j` - Move down
- `/` - Filter/search services

#### Details View (right pane)
- `↑`/`k` - Scroll up
- `↓`/`j` - Scroll down
- `Ctrl+K`/`PgUp` - Page up
- `Ctrl+J`/`PgDn` - Page down
- `g` - Go to top
- `G` - Go to bottom

## Architecture

The project is organized into clean, modular components:

```
mdns-browser/
├── cmd/mdns-browser/     # Main application entry point
├── internal/
│   ├── discovery/        # mDNS service discovery logic
│   │   ├── discover.go   # Core discovery implementation
│   │   ├── services.go   # 570+ supported service types
│   │   └── logger.go     # Custom logging configuration
│   ├── data/             # Data models and formatting
│   │   └── item.go       # Service item structure and rendering
│   └── tui/              # Terminal UI implementation
│       └── tui.go        # Bubble Tea TUI with list and viewport
```

### Key Technologies

- **[hashicorp/mdns](https://github.com/hashicorp/mdns)** - mDNS/Bonjour service discovery
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Terminal UI framework
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components (list, viewport, spinner, help)
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and layout

## Supported Services

The browser supports **570+ mDNS service types**, including:

- **Web Services**: `http`, `https`, `webdav`
- **File Sharing**: `smb`, `afpovertcp`, `nfs`, `ftp`
- **Media**: `airplay`, `raop`, `daap` (iTunes), `spotify-connect`
- **Printers**: `printer`, `ipp`
- **Remote Access**: `ssh`, `sftp-ssh`, `telnet`, `rdp`
- **Smart Home**: `homekit`, `nest`, `hue`
- **Development**: `distcc`, `git`, `svn`
- **Communication**: `sip`, `xmpp`, `irc`
- **And many more...**

See [`internal/discovery/services.go`](internal/discovery/services.go) for the complete list.

## Building

The project includes a comprehensive Makefile with multiple targets:

```bash
# Build for current platform
make build

# Run tests
make test

# Lint code (uses golangci-lint in Docker)
make lint

# Check code formatting
make format-check

# Format code
make format

# Build for all platforms (cross-compile)
make release

# Clean build artifacts
make clean

# Install dependencies
make deps
```

## Development

### Prerequisites

- Go 1.25 or later
- Docker (for running golangci-lint)
- goimports (for code formatting)

### Running Locally

```bash
# Install dependencies
make deps

# Run the application
make run

# Or run directly with go
go run cmd/mdns-browser/main.go
```

### Code Quality

Before submitting changes, ensure your code passes all checks:

```bash
# Format your code
make format

# Run linter
make lint

# Run tests
make test
```

## CI/CD

The project uses GitHub Actions for continuous integration and automated releases:

- **Lint**: Code quality checks with golangci-lint
- **Format Check**: Ensures code is properly formatted
- **Test**: Runs the test suite
- **Build**: Compiles the binary
- **Release**: Automatically creates releases with binaries for all platforms when version tags are pushed

To create a release, simply tag your commit:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with the excellent [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- mDNS discovery powered by [hashicorp/mdns](https://github.com/hashicorp/mdns)
- Service type list based on the official [IANA Service Name and Transport Protocol Port Number Registry](http://www.dns-sd.org/ServiceTypes.html)

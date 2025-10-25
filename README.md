# Sun CLI

> A powerful command-line tool built with Go

[![CI/CD](https://github.com/itsiqbal/sun-cli/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/itsiqbal/sun-cli/actions/workflows/ci-cd.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/itsiqbal/sun-cli)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/itsiqbal/sun-cli)](https://github.com/itsiqbal/sun-cli/releases/latest)
[![License](https://img.shields.io/github/license/itsiqbal/sun-cli)](LICENSE)

## Features

- 🚀 Fast and efficient Go-based CLI
- 🔒 Built with security in mind
- 📦 Easy installation with one-line installer
- 🍎 Native support for macOS (Intel & Apple Silicon)
- 🐧 Linux support (AMD64, ARM64, 386)
- 🪟 Windows support
- 🔄 Automatic updates available

## Installation

### Quick Install (macOS & Linux)

#### One-line installer (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/itsiqbal/sun-cli/main/install.sh | sudo bash
```

Or without sudo (installs to `~/.local/bin`):

```bash
curl -fsSL https://raw.githubusercontent.com/itsiqbal/sun-cli/main/install.sh | bash
```

### Manual Installation

#### macOS

**Apple Silicon (M1/M2/M3/M4):**

```bash
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-darwin-arm64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

**Intel Mac:**

```bash
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-darwin-amd64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

**Universal Binary (works on both):**

```bash
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-darwin-universal -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

#### Linux

**AMD64:**

```bash
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-linux-amd64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

**ARM64:**

```bash
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-linux-arm64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

#### Windows

**PowerShell:**

```powershell
Invoke-WebRequest -Uri "https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-windows-amd64.exe" -OutFile "sun.exe"
```

Then add the directory containing `sun.exe` to your PATH.

### Alternative Installation Methods

#### Using Go Install

```bash
go install github.com/itsiqbal/sun-cli@latest
```

#### From Source

```bash
git clone https://github.com/itsiqbal/sun-cli.git
cd sun-cli
make build
sudo make install
```

## Usage

```bash
# Show version
sun --version

# Show help
sun --help

# [Add your specific commands here]
sun [command] [flags]
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Building

```bash
# Clone the repository
git clone https://github.com/itsiqbal/sun-cli.git
cd sun-cli

# Build for current platform
make build

# Build for all platforms
make build-all

# Build for macOS only
make build-mac

# Build for Linux only
make build-linux

# Build for Windows only
make build-windows
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make coverage

# Run linters
make lint

# Run all checks
make check
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build for current platform
make build-fast    # Fast build without optimizations
make build-all     # Build for all platforms
make build-mac     # Build for macOS (Intel + Apple Silicon)
make build-linux   # Build for Linux (AMD64 + ARM64)
make build-windows # Build for Windows
make test          # Run tests with coverage
make coverage      # Generate HTML coverage report
make lint          # Run linters with autofix
make fmt           # Format code
make vet           # Run go vet
make check         # Run all checks (fmt, vet, lint, test)
make benchmark     # Run benchmarks
make clean         # Remove build artifacts
make install       # Install to $GOPATH/bin
make uninstall     # Remove from $GOPATH/bin
make release       # Create release builds
make deps          # Download dependencies
make tidy          # Tidy dependencies
make upgrade       # Upgrade all dependencies
```

## Updating

### Using the installer

```bash
curl -fsSL https://raw.githubusercontent.com/itsiqbal/sun-cli/main/install.sh | sudo bash
```

### Using Go

```bash
go install github.com/itsiqbal/sun-cli@latest
```

### Manual update

Download the latest binary from the [releases page](https://github.com/itsiqbal/sun-cli/releases/latest) and replace your existing binary.

## Uninstalling

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/sun

# If installed to ~/.local/bin
rm ~/.local/bin/sun

# If installed via go install
rm $(go env GOPATH)/bin/sun
```

## Troubleshooting

### macOS: "sun" cannot be opened because the developer cannot be verified

Remove the quarantine attribute:

```bash
xattr -d com.apple.quarantine /usr/local/bin/sun
```

### Permission denied

Use sudo for system-wide installation:

```bash
sudo curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-darwin-arm64 -o /usr/local/bin/sun
sudo chmod +x /usr/local/bin/sun
```

Or install to user directory:

```bash
mkdir -p ~/.local/bin
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-darwin-arm64 -o ~/.local/bin/sun
chmod +x ~/.local/bin/sun
export PATH="$PATH:$HOME/.local/bin"
```

### Command not found after installation

Add the installation directory to your PATH:

**For zsh (macOS default):**

```bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
source ~/.zshrc
```

**For bash:**

```bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc
```

## Verifying Downloads

All releases include checksums. To verify your download:

```bash
# Download the checksums file
curl -L https://github.com/itsiqbal/sun-cli/releases/latest/download/sun-cli_VERSION_checksums.txt -o checksums.txt

# Verify (macOS/Linux)
shasum -a 256 -c checksums.txt 2>&1 | grep sun-darwin-arm64

# Or manually
shasum -a 256 sun-darwin-arm64
```

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Linting**: Code quality checks with golangci-lint
- **Testing**: Automated tests on multiple Go versions
- **Building**: Multi-platform binary compilation
- **Security**: Vulnerability scanning with Gosec and Trivy
- **Releasing**: Automated releases with GoReleaser
- **Benchmarking**: Performance testing

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Commit Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `test:` - Test changes
- `perf:` - Performance improvements
- `refactor:` - Code refactoring

## License

[MIT License](LICENSE) - see the LICENSE file for details

## Support

- 📖 [Documentation](https://github.com/itsiqbal/sun-cli/wiki)
- 🐛 [Issue Tracker](https://github.com/itsiqbal/sun-cli/issues)
- 💬 [Discussions](https://github.com/itsiqbal/sun-cli/discussions)

## Acknowledgments

Built with ❤️ using:

- [Go](https://go.dev/)
- [GoReleaser](https://goreleaser.com/)
- [GitHub Actions](https://github.com/features/actions)

---

**[⬆ back to top](#sun-cli)**

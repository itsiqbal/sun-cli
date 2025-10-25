# Installation

## Quick Install (macOS)

### One-line installer

```bash
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/install.sh | bash
```

### Manual Installation

#### Apple Silicon (M1/M2/M3/M4)

```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-darwin-arm64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

#### Intel Mac

```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-darwin-amd64 -o /usr/local/bin/sun
chmod +x /usr/local/bin/sun
```

## Alternative Installation Methods

### Using Go Install

If you have Go installed:

```bash
go install github.com/YOUR_USERNAME/YOUR_REPO@latest
```

### Using Makefile (for developers)

Clone the repository and build locally:

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO
make build          # Build for current platform
make build-mac      # Build for both Intel and Apple Silicon
make install        # Build and install to $GOPATH/bin
```

## Verify Installation

```bash
sun --version
sun --help
```

## Updating

Simply run the installation command again to update to the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/install.sh | bash
```

Or with Go:

```bash
go install github.com/YOUR_USERNAME/YOUR_REPO@latest
```

## Uninstalling

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/sun

# If installed to ~/.local/bin
rm ~/.local/bin/sun

# If installed via go install
rm $(go env GOPATH)/bin/sun
```

## Building from Source

### Quick Build

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO
make build
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build for current platform
make build-mac     # Build for both macOS architectures
make build-all     # Build for all platforms
make test          # Run tests with coverage
make lint          # Run linters
make check         # Run all checks (fmt, vet, lint, test)
make clean         # Clean build artifacts
make release       # Create release builds
```

## Checksums

To verify your download:

```bash
# Download checksums file
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun_VERSION_checksums.txt -o checksums.txt

# Verify (macOS)
shasum -a 256 -c checksums.txt 2>&1 | grep sun-darwin
```

## Troubleshooting

### "sun" cannot be opened because the developer cannot be verified

If you see this error on macOS:

```bash
# Remove the quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/sun
```

### Permission denied

If you get permission errors:

```bash
# Use sudo for system-wide installation
sudo curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-darwin-arm64 -o /usr/local/bin/sun
sudo chmod +x /usr/local/bin/sun

# Or install to user directory
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-darwin-arm64 -o ~/.local/bin/sun
chmod +x ~/.local/bin/sun
```

### Command not found after installation

Add the installation directory to your PATH:

```bash
# For ~/.local/bin (add to ~/.zshrc or ~/.bashrc)
export PATH="$PATH:$HOME/.local/bin"

# Then reload your shell
source ~/.zshrc  # or source ~/.bashrc
```

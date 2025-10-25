# Installation

## Quick Install (macOS & Linux)

### One-line installer
```bash
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/YOUR_REPO/main/install.sh | bash
```

### Manual Installation

#### macOS

**Apple Silicon (M1/M2/M3/M4):**
```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli-darwin-arm64 -o /usr/local/bin/sun-cli
chmod +x /usr/local/bin/sun-cli
```

**Intel Mac:**
```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli-darwin-amd64 -o /usr/local/bin/sun-cli
chmod +x /usr/local/bin/sun-cli
```

#### Linux

**ARM64:**
```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli-linux-arm64 -o /usr/local/bin/sun-cli
chmod +x /usr/local/bin/sun-cli
```

**AMD64:**
```bash
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli-linux-amd64 -o /usr/local/bin/sun-cli
chmod +x /usr/local/bin/sun-cli
```

#### Windows

Download the latest release:
```powershell
# PowerShell
Invoke-WebRequest -Uri "https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli-windows-amd64.exe" -OutFile "sun-cli.exe"
```

Then add the directory containing `sun-cli.exe` to your PATH.

## Alternative: Go Install

If you have Go installed:
```bash
go install github.com/YOUR_USERNAME/YOUR_REPO@latest
```

## Verify Installation

```bash
sun-cli --version
sun-cli --help
```

## Updating

Simply run the installation command again to update to the latest version.

## Uninstalling

```bash
# macOS & Linux
sudo rm /usr/local/bin/sun-cli

# Or if installed to ~/.local/bin
rm ~/.local/bin/sun-cli
```

## Building from Source

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO
go build -o sun-cli .
sudo mv sun-cli /usr/local/bin/
```

## Checksums

To verify your download:
```bash
# Download checksums
curl -L https://github.com/YOUR_USERNAME/YOUR_REPO/releases/latest/download/sun-cli_VERSION_checksums.txt -o checksums.txt

# Verify (macOS/Linux)
sha256sum -c checksums.txt 2>&1 | grep sun-cli-darwin-arm64
```
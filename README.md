<div align="center">

# gen

![Go Version](https://img.shields.io/github/go-mod/go-version/aminshahid573/gen?style=for-the-badge&labelColor=101418&color=9ccbfb)
![Latest Release](https://img.shields.io/github/v/release/aminshahid573/gen?style=for-the-badge&labelColor=101418&color=b9c8da)
![GitHub last commit](https://img.shields.io/github/last-commit/aminshahid573/gen?style=for-the-badge&labelColor=101418&color=d3bfe6)
![GitHub repo size](https://img.shields.io/github/repo-size/aminshahid573/gen?style=for-the-badge&labelColor=101418&color=96f1f1)

</div>

A collection of small dev utilities packed into a single CLI — things you'd normally google or visit a website for. UUIDs, passwords, QR codes, OTPs, lorem ipsum, tokens and more, all from your terminal.

## What's in it

- **id** — UUIDs and other unique identifiers
- **str** — string manipulation and transforms
- **token** — secure tokens for APIs and auth
- **qr** — render QR codes right in your terminal
- **otp** — TOTP/HOTP codes for testing 2FA flows
- **lorem** — placeholder text for mockups and layouts
- **time** — timestamps, timezones, formatting
- **pin** — numeric PINs for auth systems
- **pass** — strong random passwords
- More being added as needed

## Installation

**Homebrew** (macOS / Linux)
```bash
brew install aminshahid573/tap/gen
```

**Scoop** (Windows)
```bash
scoop bucket add aminshahid573 https://github.com/aminshahid573/scoop-bucket
scoop install gen
```

**apt** (Debian / Ubuntu)
```bash
wget https://github.com/aminshahid573/gen/releases/latest/download/gen_1.0.6_linux_amd64.deb
sudo dpkg -i gen_1.0.6_linux_amd64.deb
```

**rpm** (Fedora / RHEL)
```bash
rpm -i https://github.com/aminshahid573/gen/releases/download/v1.0.6/gen_1.0.6_linux_amd64.rpm
```

**apk** (Alpine)
```bash
wget https://github.com/aminshahid573/gen/releases/latest/download/gen_1.0.6_linux_amd64.apk
apk add --allow-untrusted gen_1.0.6_linux_amd64.apk
```

**curl / wget** (any OS)
```bash
curl -fsSL https://raw.githubusercontent.com/aminshahid573/gen/main/install.sh | sh
```

**go install**
```bash
go install github.com/aminshahid573/gen@latest
```

**Build from source**
```bash
git clone https://github.com/aminshahid573/gen.git
cd gen
make build
```

## Usage

```bash
gen --help           # see all commands
gen qr --help        # options for a specific command
```

### Quick examples

```bash
# QR code for a URL
gen qr "https://example.com"

# TOTP code for 2FA testing
gen otp --secret JBSWY3DPEHPK3PXP

# 16-character password
gen pass --length 16

# UUID v4
gen id --type uuid4
```

## Development

```bash
go mod tidy       # install dependencies
make test         # run tests
make fmt          # format code
make build        # build binary
make clean        # clean artifacts
```

## Contributing

Open an issue if something's broken or you have an idea for a new command. PRs are welcome.

1. Fork the repo
2. Create a branch (`git checkout -b feature/some-command`)
3. Commit and push
4. Open a PR against `main`

## License

MIT — see [LICENSE](LICENSE) for details.

Built with [Cobra](https://github.com/spf13/cobra).

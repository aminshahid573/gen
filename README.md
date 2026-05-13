# gen

A Swiss Army knife CLI tool for developers, packed with useful utilities for daily tasks.

## What it does

gen is a collection of command-line tools that make common development tasks quicker and less painful. Instead of remembering obscure commands or visiting websites for simple utilities, you've got them right in your terminal.

## Features

- **ID Generation** - Create UUIDs and other unique identifiers
- **String Utilities** - Manipulate and transform text in various ways
- **Token Generation** - Generate secure tokens for APIs, authentication, and more
- **QR Code Generator** - Create QR codes directly in your terminal with various error correction levels
- **OTP Generator** - Generate TOTP and HOTP codes for 2FA testing and validation
- **Lorem Ipsum** - Generate placeholder text when you need to mock up designs or test layouts
- **Time Utilities** - Work with timestamps, time zones, and time formatting
- **PIN Generation** - Create secure numeric PINs for authentication systems
- **Password Generation** - Generate strong, memorable passwords
- And more utilities being added regularly...

## Installation

```bash
# Install via go install
go install github.com/aminshahid573/gen@latest

# Or build from source
git clone https://github.com/aminshahid573/gen.git
cd gen
make build
```

## Usage

Each utility has its own help documentation. Try:

```bash
gen --help              # See all available commands
gen qr --help           # QR code generator options
gen otp --help          # OTP generator options
# ... and so on for each command
```

### Examples

Generate a QR code for a website:
```bash
gen qr "https://example.com"
```

Generate a TOTP code for testing 2FA:
```bash
gen otp --secret JBSWY3DPEHPK3PXP
```

Create a secure password:
```bash
gen pass --length 16
```

Generate a UUID:
```bash
gen id --type uuid4
```

## Development

gen is built with Go and uses the Cobra library for command structure.

```bash
# Install dependencies
go mod tidy

# Run tests
make test

# Format code
make fmt

# Build the binary
make build

# Clean build artifacts
make clean
```

## Contributing

Found a bug or have an idea for a new utility? Feel free to open an issue or submit a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-tool`)
3. Commit your changes (`git commit -am 'Add amazing tool'`)
4. Push to the branch (`git push origin feature/amazing-tool`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - For powerful CLI command handling
- Various other open-source libraries that make development easier

---

Created with ❤️ by developers, for developers. If you find this useful, a star on GitHub is always appreciated!

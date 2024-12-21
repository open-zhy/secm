# secm - Secure Secret Manager

A command-line tool for securely managing secrets with encryption and metadata support.

## Features

- Secure storage of secrets using hybrid encryption (RSA + AES)
- YAML-based secret storage with metadata
- Support for secret tags and categorization
- Cross-platform support (Linux, macOS, Windows)

## Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/open-zhy/secm.git
cd secm
```

2. Build for your platform:
```bash
make
```

Or build for a specific platform:
```bash
make build-platform PLATFORM=darwin ARCH=arm64
```

Build for all platforms:
```bash
make build-all
```

## Usage

### Initialize Workspace

Before using secm, initialize the workspace:

```bash
secm init
```

This creates the `.secm` directory in your home folder and generates an RSA identity key.

### Create a Secret

Create a new secret from a file with metadata:

```bash
secm create secret.txt -n "API Key" -d "Production API key" -t "api,prod" --type "api-key"
```

Options:
- `-n, --name`: Name of the secret (required)
- `-d, --description`: Description of the secret
- `-t, --type`: Type of secret (e.g., api-key, certificate)
- `--tags`: Comma-separated list of tags
- `-f, --format`: Format of the secret (text, json, binary)

### List Secrets

List all stored secrets:

```bash
secm list
```

Show additional information:
```bash
secm list -t    # Show tags
secm list -d    # Show descriptions
```

### Get a Secret

Retrieve a secret by its ID:

```bash
secm get <secret-id>                    # Output to stdout
secm get <secret-id> -o output.txt      # Save to file
secm get <secret-id> -m                 # Show metadata
secm get <secret-id> -q                 # Quiet mode (only output value)
```

## Building from Source

Requirements:
- Go 1.16 or later
- Make

Available make commands:
- `make`: Build for current platform
- `make build-all`: Build for all platforms
- `make build-platform PLATFORM=darwin ARCH=arm64`: Build for specific platform
- `make clean`: Clean build directory
- `make test`: Run tests
- `make fmt`: Format code
- `make install`: Install locally

## Security

- Uses hybrid encryption (RSA for key exchange, AES for data)
- Secure file permissions (0600 for keys, 0700 for directories)
- Unique hash-based IDs for secrets
- Base64 encoded encrypted data in YAML storage

## Todo

- [x] Basics of secrets management: workspace initialization, create secret, list and unfold secret
- [ ] Add `--profile` option on root level, default to `~/.secm`: this should enable multiple instances or easily resurrect from an existing profile
- [ ] Support `ed25519` key and more
- [ ] Enable transfer to another identity: `secm transfer <publicKey>`: it will just create a copy in the workspace of the same secret, only recipient can read the secret
- [ ] After transfer, enable p2p direct transfer (preferrable implemented as plugin, not apart of the core util)

## License

MIT License

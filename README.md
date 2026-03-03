# echo-cli

A CLI tool for bootstrapping Echo web server projects. Not affiliated with the official LabStack Echo framework.

## Installation

### Go Install

```bash
go install github.com/cjodo/echo-cli/cmd/echo-cli
```

### From Source

```bash
git clone https://github.com/cjodo/echo-cli.git
cd echo-cli
go build -o echo-cli ./cmd/echo-cli.
```

### Verify Installation

```bash
echo-cli --help
```

## Usage

### Cookbook

Browse and pull ready-made examples from the official Echo cookbook.

```bash
echo-cli cookbook list                    # List available recipes
echo-cli cookbook get <recipe-name>        # Download a recipe to current directory
```

Flags:
- `--refresh` - Force refresh cache, bypasses cached API responses and downloaded files

Examples:

```bash
# List all available cookbook recipes
echo-cli cookbook list

# Download a recipe (cached after first run)
echo-cli cookbook get file-download

# Force refresh from GitHub
echo-cli cookbook get file-download --refresh

# Force refresh recipe list
echo-cli cookbook list --refresh
```

Available recipes:
- auto-tls
- casbin
- cors
- crud
- csrf
- embed
- file-download
- file-upload
- graceful-shutdown
- hello-world
- http2-server-push
- http2
- jsonp
- jwt
- load-balancing
- middleware
- prometheus
- reverse-proxy
- sse
- streaming-response
- subdomain
- timeout
- websocket

### Docs

Serve the Echo framework documentation locally for offline access.

```bash
echo-cli docs              # Serve docs on default port 8080
echo-cli docs -p 3000      # Serve docs on port 3000
echo-cli docs --refresh    # Force refresh docs cache
```

Flags:
- `-p, --port string`   Port to serve docs on (default "8080")
- `--refresh`           Force refresh cache, bypasses cached docs

Examples:

```bash
# Serve docs locally
echo-cli docs

# Serve on custom port
echo-cli docs -p 8081

# Refresh the docs cache
echo-cli docs --refresh
```

## Cache

The cookbook and docs commands cache GitHub API responses and downloaded files to speed up subsequent requests.

Cache location: `~/.echo-cli/cache/`

- `api/` - GitHub API responses (directory listings)
- `files/` - Downloaded recipe and docs files

Cache TTL: 24 hours (default)

Use `--refresh` flag to bypass the cache and fetch fresh data from GitHub.

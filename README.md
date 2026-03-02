# echo-cli

A CLI tool for bootstrapping Echo web server projects. Not affiliated with the official LabStack Echo framework.

## Installation

### From Source

```bash
git clone https://github.com/cjodo/echo-cli.git
cd echo-cli
go build -o echo ./echo/
```

Or install globally:

```bash
go install ./echo/
```

### Verify Installation

```bash
echo --help
```

## Usage

### Create a New Project

```bash
echo new <project-name> [module-name]
```

Creates a new Echo project in a directory with the specified name.

Arguments:
- `project-name` (required): Name of the project directory to create
- `module-name` (optional): Go module name (defaults to project-name)

### Use a Template

```bash
echo new my-api -t <template>
```

Available templates:

- `auto-tls` - Automatic TLS with Let's Encrypt
- `cors` - CORS middleware configuration
- `crud` - Basic CRUD API endpoints
- `embed-resources` - Embedded static assets
- `file-download` - File download handler
- `file-upload` - File upload handler
- `graceful-shutdown` - Graceful server shutdown
- `hello-world` - Basic Echo server (default)
- `http2-server` - HTTP/2 server
- `http2-server-push` - HTTP/2 server push
- `jsonp` - JSONP response handling
- `jwt` - JWT authentication
- `load-balancing` - Load balancing with round-robin
- `middleware` - Custom middleware example
- `reverse-proxy` - Reverse proxy setup
- `streaming-response` - Streaming JSON response
- `subdomain` - Virtual host subdomain routing
- `timeout` - Request timeout handling
- `websocket` - WebSocket support

> [!NOTE]
> All templates were pulled from the [https://echo.labstack.com/docs/category/cookbook](https://echo.labstack.com/docs/category/cookbook) examples

### Examples

Create a default hello-world project:

```bash
echo new my-project
```

Create a project with JWT auth:

```bash
echo new api -t jwt
```

Create a project with a custom module name:

```bash
echo new my-project github.com/myuser/my-project
```

## Project Structure

After running `echo new my-project`, the generated project will have:

```
my-project/
├── go.mod
├── server.go
└── [additional template files]
```

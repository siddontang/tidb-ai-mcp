# TiDB AI MCP Server

A Model Control Protocol (MCP) server that allows Cursor to ask questions to TiDB AI.

## Features

- Server-Sent Events (SSE) transport for MCP protocol
- Integration with TiDB AI for answering questions
- Simple and lightweight implementation

## Installation

```bash
go get github.com/siddontang/tidb-ai-mcp
```

## Usage

The server supports a command-line argument to configure the port:

```bash
# Show help
go run main.go -h

# Available options:
# -port int
#   Port to listen on for SSE transport (default 8080)
```

### Running the Server

```bash
# Start the MCP server with default settings
go run main.go

# Specify a custom port
go run main.go -port 9000
```

The server will listen on the specified port with the endpoint `/sse`. By default, it listens on port 8080.

## Using with Cursor

To use this MCP server with Cursor:

1. Start the server: `go run main.go`
2. In Cursor, configure the MCP server URL to `http://localhost:8080/sse` (or your custom port)
3. You can now ask questions to TiDB AI directly from Cursor

## License

MIT 
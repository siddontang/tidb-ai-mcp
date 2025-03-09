# TiDB AI MCP Server

A Model Control Protocol (MCP) server that allows Cursor to ask questions to TiDB AI.

## Features

- Allows Cursor to ask questions to TiDB AI
- Uses stdio transport for communication with Cursor
- Simple command-line interface

## Installation

```bash
go get github.com/siddontang/tidb-ai-mcp
```

## Usage

```bash
# Run the server with stdio transport
./tidb-ai-mcp
```

## Building

```bash
go build -o tidb-ai-mcp
```

## Using with Cursor

To use this MCP server with Cursor:

1. Start the server: `go run main.go`
2. In Cursor, configure the MCP server URL to `http://localhost:8080/sse` (or your custom port)
3. You can now ask questions to TiDB AI directly from Cursor

## License

MIT 
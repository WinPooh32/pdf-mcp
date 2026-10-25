# pdf-mcp

MCP server for reading PDFs: the `read` tool renders a given PDF page to PNG and returns it to the LLM as image content.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-http` | (empty) | streamable-HTTP address; empty = stdio |
| `-dpi` | `150` | rendering resolution, dots per inch |

## Build and run

```sh
go install github.com/WinPooh32/pdf-mcp@latest
```

### stdio (default)

```sh
$ ~/go/bin/pdf-mcp
```

Connecting from an MCP client (example, Claude Desktop):

```json
{
  "mcpServers": {
    "pdf": {
      "command": "/path/to/pdf-mcp",
      "args": ["-dpi", "150"]
    }
  }
}
```

### streamable-HTTP

```sh
./pdf-mcp -http :9000
```

The client connects to `http://127.0.0.1:9000/mcp` (stateless, JSON responses).

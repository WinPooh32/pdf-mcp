package mcputil

import "github.com/modelcontextprotocol/go-sdk/mcp"

func TextResult(txt string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text:        txt,
				Meta:        nil,
				Annotations: nil,
			},
		},
		Meta:              nil,
		StructuredContent: nil,
		IsError:           false,
	}
}

func ErrorResult(txt string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text:        txt,
				Meta:        nil,
				Annotations: nil,
			},
		},
		Meta:              nil,
		StructuredContent: nil,
		IsError:           true,
	}
}

// ImageResult wraps raw image bytes; the SDK base64-encodes them on the wire.
func ImageResult(data []byte) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.ImageContent{
				Data:        data,
				MIMEType:    "image/png",
				Meta:        nil,
				Annotations: nil,
			},
		},
		Meta:              nil,
		StructuredContent: nil,
		IsError:           false,
	}
}

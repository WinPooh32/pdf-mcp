package tools

import (
	"fmt"

	"github.com/WinPooh32/pdf-mcp/tools/read"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Kit struct {
	server *mcp.Server
	tools  []mcp.Tool
}

func (kit *Kit) Server() *mcp.Server {
	return kit.server
}

func (kit *Kit) Add(t mcp.Tool) error {
	if t.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	kit.tools = append(kit.tools, t)

	return nil
}

func NewDefault(config *ServerConfig) *Kit {
	kit := New()

	// Set up read config
	read.SetConfig(&read.Config{
		Dpi: config.Dpi,
	})

	mcpToolRead := &mcp.Tool{
		Name:         "read",
		Description:  "Read a PDF page and return it as a PNG image",
		Meta:         nil,
		Annotations:  nil,
		InputSchema:  read.Schema(),
		OutputSchema: nil,
		Title:        "",
		Icons:        nil,
	}
	mcp.AddTool(kit.server, mcpToolRead, read.Read)

	return kit
}

func New() *Kit {
	return &Kit{
		server: mcp.NewServer(&mcp.Implementation{Name: "pdf-mcp", Version: "0.0.1", Title: "", WebsiteURL: "", Icons: nil}, nil),
		tools:  []mcp.Tool{},
	}
}

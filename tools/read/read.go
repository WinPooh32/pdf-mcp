package read

import (
	"context"
	"fmt"
	"sync"

	"github.com/WinPooh32/pdf-mcp/pkg/mcputil"
	"github.com/gen2brain/go-fitz"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config struct {
	Dpi float64
}

var (
	config *Config
	mut    sync.RWMutex
)

func SetConfig(c *Config) {
	mut.Lock()
	defer mut.Unlock()
	config = c
}

func getConfig() *Config {
	mut.RLock()
	defer mut.RUnlock()
	return config
}

type ReadInput struct {
	Path string `json:"path"`
	Page int    `json:"page"`
}

func Read(_ context.Context, _ *mcp.CallToolRequest, input ReadInput) (*mcp.CallToolResult, any, error) {
	if input.Path == "" {
		return mcputil.ErrorResult("path is required"), nil, nil
	}

	if input.Page < 1 {
		return mcputil.ErrorResult(fmt.Sprintf("page must be at least 1, got %d", input.Page)), nil, nil
	}

	dpi := getConfig().Dpi

	doc, err := fitz.New(input.Path)
	if err != nil {
		return mcputil.ErrorResult(fmt.Sprintf("open pdf: %v", err)), nil, nil
	}
	defer func() { _ = doc.Close() }()

	if input.Page > doc.NumPage() {
		return mcputil.ErrorResult(fmt.Sprintf("page %d out of range (document has %d pages)", input.Page, doc.NumPage())), nil, nil
	}

	png, err := doc.ImagePNG(input.Page-1, dpi)
	if err != nil {
		return mcputil.ErrorResult(fmt.Sprintf("render page: %v", err)), nil, nil
	}

	return mcputil.ImageResult(png), nil, nil
}

// Schema returns the JSON schema for ReadInput.
func Schema() any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"description": "Path to the PDF file",
				"type":        "string",
			},
			"page": map[string]any{
				"description": "Page number to render (1-based)",
				"type":        "integer",
			},
		},
		"required": []string{"path", "page"},
	}
}

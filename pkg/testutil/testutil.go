package testutil

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CreateSessionPair creates a client-server pair via in-memory transports.
func CreateSessionPair(t *testing.T) (*mcp.Client, *mcp.Server, *mcp.ClientSession, *mcp.ServerSession) {
	t.Helper()

	t1, t2 := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "pdf-mcp", Version: "0.0.1", Title: "", WebsiteURL: "", Icons: nil}, nil)
	serverSession, err := server.Connect(context.Background(), t1, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = serverSession.Close() })

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1", Title: "", WebsiteURL: "", Icons: nil}, nil)
	clientSession, err := client.Connect(context.Background(), t2, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = clientSession.Close() })

	return client, server, clientSession, serverSession
}

// AssertToolSuccess asserts a successful tool call.
func AssertToolSuccess(t *testing.T, result *mcp.CallToolResult) *mcp.CallToolResult {
	t.Helper()

	assert.False(t, result.IsError, "expected tool call to succeed")
	assert.Len(t, result.Content, 1, "expected single content item")

	return result
}

// AssertToolError asserts a failed tool call with a message containing expectedErrorSubstring.
func AssertToolError(t *testing.T, result *mcp.CallToolResult, expectedErrorSubstring string) {
	t.Helper()

	assert.True(t, result.IsError, "expected tool call to fail")
	assert.Len(t, result.Content, 1, "expected single content item")
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "expected TextContent")
	assert.Contains(t, textContent.Text, expectedErrorSubstring, "expected error to contain %q", expectedErrorSubstring)
}

// AddToolToServer adds a tool to the server.
func AddToolToServer(server *mcp.Server, name string, description string, inputSchema any, handler mcp.ToolHandler) {
	mcpTool := &mcp.Tool{
		Name:         name,
		Description:  description,
		InputSchema:  inputSchema,
		Meta:         nil,
		Annotations:  nil,
		OutputSchema: nil,
		Title:        "",
		Icons:        nil,
	}
	server.AddTool(mcpTool, handler)
}

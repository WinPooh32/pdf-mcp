package read_test

import (
	"context"
	"encoding/json"
	"image/png"
	"bytes"
	"testing"

	"github.com/WinPooh32/pdf-mcp/pkg/mcputil"
	"github.com/WinPooh32/pdf-mcp/pkg/testutil"
	"github.com/WinPooh32/pdf-mcp/tools/read"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testPDF is an external 5-page A4 PDF used as the fixture.
const testPDF = "/run/media/system/Data/Specification - Model Context Protocol.pdf"

// addReadTool registers the read tool on the server and returns a caller.
func addReadTool(t *testing.T) (clientSession *mcp.ClientSession, call func(path string, page int) *mcp.CallToolResult) {
	t.Helper()

	_, server, clientSession, serverSession := testutil.CreateSessionPair(t)
	t.Cleanup(func() {
		_ = clientSession.Close()
		_ = serverSession.Close()
	})

	handler := func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var input read.ReadInput
		if err := json.Unmarshal(request.Params.Arguments, &input); err != nil {
			return mcputil.ErrorResult("invalid arguments"), nil
		}

		result, _, err := read.Read(ctx, request, input)
		if err != nil {
			return result, nil
		}
		return result, nil
	}

	testutil.AddToolToServer(server, "read", "Render a PDF page as a PNG image", read.Schema(), handler)

	return clientSession, func(path string, page int) *mcp.CallToolResult {
		t.Helper()

		result, err := clientSession.CallTool(context.Background(), &mcp.CallToolParams{
			Name:      "read",
			Arguments: map[string]any{"path": path, "page": page},
			Meta:      nil,
		})
		require.NoError(t, err)

		return result
	}
}

func assertPNG(t *testing.T, result *mcp.CallToolResult) {
	t.Helper()

	testutil.AssertToolSuccess(t, result)

	imageContent, ok := result.Content[0].(*mcp.ImageContent)
	require.True(t, ok, "expected ImageContent")
	assert.Equal(t, "image/png", imageContent.MIMEType)

	img, err := png.Decode(bytes.NewReader(imageContent.Data))
	require.NoError(t, err, "expected valid PNG")
	assert.NotEmpty(t, img.Bounds().Size())
}

func TestRead(t *testing.T) {
	// Not parallel - uses global config
	read.SetConfig(&read.Config{Dpi: 150})

	clientSession, call := addReadTool(t)
	assertPNG(t, call(testPDF, 1))
	_ = clientSession
}

func TestRead_LastPage(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	assertPNG(t, call(testPDF, 5))
}

func TestRead_OutOfRange(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	testutil.AssertToolError(t, call(testPDF, 6), "page 6 out of range (document has 5 pages)")
}

func TestRead_ZeroPage(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	testutil.AssertToolError(t, call(testPDF, 0), "page must be at least 1")
}

func TestRead_NegativePage(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	testutil.AssertToolError(t, call(testPDF, -1), "page must be at least 1")
}

func TestRead_EmptyPath(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	testutil.AssertToolError(t, call("", 1), "path is required")
}

func TestRead_FileNotFound(t *testing.T) {
	read.SetConfig(&read.Config{Dpi: 150})

	_, call := addReadTool(t)
	testutil.AssertToolError(t, call("/nope/definitely-missing.pdf", 1), "open pdf")
}

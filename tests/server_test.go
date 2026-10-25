package tests

import (
	"context"
	"os"
	"testing"

	"github.com/WinPooh32/pdf-mcp/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerRunStdio(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := tools.NewServer()
	kit := tools.New()

	err := srv.Run(ctx, &tools.ServerConfig{HTTPAddr: "", StdioMode: true, HTTPMode: false}, kit)
	require.NoError(t, err)
}

func TestNewDefaultKit(t *testing.T) {
	config := &tools.ServerConfig{Dpi: 150}

	kit := tools.NewDefault(config)
	assert.NotNil(t, kit, "NewDefault should return a non-nil kit")
}

func TestNewEmptyKit(t *testing.T) {
	kit := tools.New()
	assert.NotNil(t, kit, "New should return a non-nil kit")
}

// parseArgs runs ParseArgsConfig with the given CLI args (test-binary flags stripped).
func parseArgs(t *testing.T, args ...string) (*tools.ServerConfig, error) {
	t.Helper()

	origArgs := os.Args
	t.Cleanup(func() { os.Args = origArgs })
	os.Args = append([]string{"pdf-mcp"}, args...)

	return tools.NewServer().ParseArgsConfig()
}

func TestParseArgsConfig(t *testing.T) {
	config, err := parseArgs(t)
	require.NoError(t, err)

	assert.True(t, config.StdioMode)
	assert.False(t, config.HTTPMode)
	assert.Empty(t, config.HTTPAddr)
	assert.Equal(t, 150.0, config.Dpi)
}

func TestParseArgsConfig_HTTP(t *testing.T) {
	config, err := parseArgs(t, "-http", ":9000", "-dpi", "300")
	require.NoError(t, err)

	assert.False(t, config.StdioMode)
	assert.True(t, config.HTTPMode)
	assert.Equal(t, ":9000", config.HTTPAddr)
	assert.Equal(t, 300.0, config.Dpi)
}

func TestParseArgsConfig_InvalidDpi(t *testing.T) {
	for _, dpi := range []string{"0", "-5"} {
		_, err := parseArgs(t, "-dpi", dpi)
		require.Error(t, err, "dpi %s should be invalid", dpi)
	}
}

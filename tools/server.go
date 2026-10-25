package tools

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpShutdownTimeout   = 15 * time.Second
)

type ServerConfig struct {
	HTTPAddr  string
	StdioMode bool
	HTTPMode  bool
	Dpi       float64
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (srv *Server) ParseArgsConfig() (*ServerConfig, error) {
	// New FlagSet per call: the global CommandLine cannot be re-parsed in tests.
	fs := flag.NewFlagSet("pdf-mcp", flag.ContinueOnError)
	httpAddr := fs.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
	dpi := fs.Float64("dpi", 150, "render resolution in dots per inch")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if *dpi <= 0 {
		return nil, fmt.Errorf("dpi must be positive, got %f", *dpi)
	}

	return &ServerConfig{
		HTTPAddr:  *httpAddr,
		StdioMode: *httpAddr == "",
		HTTPMode:  *httpAddr != "",
		Dpi:       *dpi,
	}, nil
}

func (srv *Server) Run(ctx context.Context, c *ServerConfig, kit *Kit) error {
	switch {
	case c.HTTPMode:
		return srv.runHTTP(ctx, c.HTTPAddr, kit)
	case c.StdioMode:
		fallthrough
	default:
		return srv.runStdio(ctx, kit)
	}
}

func (srv *Server) runStdio(ctx context.Context, kit *Kit) error {
	err := kit.server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		return fmt.Errorf("stdio server run: %w", err)
	}

	return nil
}

func (srv *Server) runHTTP(ctx context.Context, addr string, kit *Kit) error {
	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return kit.server
	}, &mcp.StreamableHTTPOptions{
		Stateless:                  true,
		JSONResponse:               true,
		Logger:                     nil,
		EventStore:                 nil,
		SessionTimeout:             0,
		DisableLocalhostProtection: false,
	})

	server := &http.Server{
		Addr:                         addr,
		Handler:                      handler,
		ReadHeaderTimeout:            httpReadHeaderTimeout,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
		DisableGeneralOptionsHandler: false,
	}

	ctx, cancel := context.WithCancel(ctx)

	var serverErr error

	go func() {
		defer cancel()

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr = fmt.Errorf("serve http: %w", err)
		}
	}()

	<-ctx.Done()

	if serverErr != nil {
		return serverErr
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), httpShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	return nil
}

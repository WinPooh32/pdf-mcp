package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/WinPooh32/pdf-mcp/tools"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := tools.NewServer()

	c, err := srv.ParseArgsConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config: %s\n", err)
		os.Exit(1)
	}

	if err := srv.Run(ctx, c, tools.NewDefault(c)); err != nil {
		fmt.Fprintf(os.Stderr, "run: %s\n", err)
		os.Exit(1)
	}
}

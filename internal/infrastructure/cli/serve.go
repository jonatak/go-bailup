package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/jonatak/go-bailup/internal/application"
)

type Server struct{}

func (s *Server) Run(ctx context.Context, server application.Server) error {
	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return nil
}

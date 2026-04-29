package cli

import (
	"fmt"

	"github.com/jonatak/go-bailup/internal/version"
)

type VersionCommand struct {
}

func (v *VersionCommand) Run() error {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Commit SHA: %s\n", version.CommitSHA)
	fmt.Printf("Build Time: %s\n", version.BuildTime)
	return nil
}

package app

import (
	"fmt"
	"os"

	"github.com/jonatak/go-bailup/internal/bailup"
)

type AppContext struct {
	BailUp *bailup.Bailup
}

func NewApp() (*AppContext, error) {
	bailupEmail := os.Getenv("BAILUP_EMAIL")
	bailupPassword := os.Getenv("BAILUP_PASS")
	bailupRegulation := os.Getenv("BAILUP_REGULATION")

	if bailupEmail == "" || bailupPassword == "" || bailupRegulation == "" {
		return nil, InitError
	}

	bailup := bailup.NewBailup(bailupEmail, bailupPassword, bailupRegulation)
	err := bailup.Connect()
	if err != nil {
		return nil, fmt.Errorf("an error occured: %v\n", err)
	}
	return &AppContext{
		BailUp: bailup,
	}, nil
}

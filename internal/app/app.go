package app

import (
	"fmt"
	"os"

	"github.com/jonatak/go-bailup/internal/application"
	"github.com/jonatak/go-bailup/internal/infrastructure/bailup"
)

type AppContext struct {
	HVACService *application.HVACService
}

func NewApp() (*AppContext, error) {
	bailupEmail := os.Getenv("BAILUP_EMAIL")
	bailupPassword := os.Getenv("BAILUP_PASS")
	bailupRegulation := os.Getenv("BAILUP_REGULATION")

	if bailupEmail == "" || bailupPassword == "" || bailupRegulation == "" {
		return nil, InitError
	}

	gateway := bailup.NewGateway(bailupEmail, bailupPassword, bailupRegulation)
	err := gateway.Connect()
	if err != nil {
		return nil, fmt.Errorf("connect HVAC gateway: %w", err)
	}
	return &AppContext{
		HVACService: application.NewHVACService(gateway),
	}, nil
}

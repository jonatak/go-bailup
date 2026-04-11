package main

import (
	"fmt"
	"os"

	"github.com/jonatak/go-bailup/internal/bailup"
)

func main() {
	fmt.Println("Welcome to bailup.")

	bailupEmail := os.Getenv("BAILUP_EMAIL")
	bailupPassword := os.Getenv("BAILUP_PASS")
	bailupRegulation := os.Getenv("BAILUP_REGULATION")

	if bailupEmail == "" || bailupPassword == "" || bailupRegulation == "" {
		fmt.Println("env var BAILUP_EMAIL, BAILUP_PASS, BAILUP_REGULATION need to be set")
		return
	}

	bailup := bailup.NewBailup(bailupEmail, bailupPassword, bailupRegulation)
	err := bailup.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

	if !bailup.IsConnected() {
		fmt.Println("Disconnected")
	}

	state, err := bailup.GetState()
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured: %v\n", err)
		return
	}

	fmt.Println(state)
}

package main

import (
	"fmt"

	"github.com/ramZenit/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading the config file:", err)
		return
	}

	fmt.Printf("main: %v\n", cfg)

}

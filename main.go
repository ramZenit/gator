package main

import (
	"errors"
	"fmt"
	"os"

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
	appState := state{cfg: &cfg}

	cmdHandlers := commands{make(map[string]func(*state, command) error)}

	cmdHandlers.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err := cmdHandlers.run(&appState, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no <username> target for login command")
	}
	username := cmd.args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("Access granted to:", username)
	return nil
}

// run a given command with the provided state if exists
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	return handler(s, cmd)
}

// registers a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

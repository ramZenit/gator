package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
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

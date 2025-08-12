package main

import (
	"context"
	"fmt"

	"github.com/ramZenit/gator/internal/database"
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

func loadCMDs() commands {
	cmdHandlers := commands{make(map[string]func(*state, command) error)}

	cmdHandlers.register("login", handlerLogin)
	cmdHandlers.register("register", handlerRegister)
	cmdHandlers.register("reset", handlerReset)
	cmdHandlers.register("users", handlerListUsers)
	cmdHandlers.register("agg", handlerAggregator)
	cmdHandlers.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmdHandlers.register("feeds", handlerListFeeds)
	cmdHandlers.register("follow", middlewareLoggedIn(handlerCreateFollow))
	cmdHandlers.register("following", middlewareLoggedIn(handlerFollowsPerUser))
	cmdHandlers.register("unfollow", middlewareLoggedIn(handlerUnfollowPerUser))
	cmdHandlers.register("browse", middlewareLoggedIn(handlerBrowse))
	return cmdHandlers

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("unable to retrieve user info: %w", err)
		}
		return handler(s, cmd, user)
	}
}

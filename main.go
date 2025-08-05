package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ramZenit/gator/internal/config"
	"github.com/ramZenit/gator/internal/database"
)

type state struct {
	db  *database.Queries
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

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println("Error opening the db:", err)
		return
	}
	dbQueries := database.New(db)

	appState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmdHandlers := commands{make(map[string]func(*state, command) error)}

	cmdHandlers.register("login", handlerLogin)
	cmdHandlers.register("register", handlerRegister)
	cmdHandlers.register("reset", handlerReset)
	cmdHandlers.register("users", handlerUsers)

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
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("username %s not found", username)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("Access granted to user:", username)
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no <user> target for register command")
	}
	user := cmd.args[0]
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      user,
	}
	res, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	s.cfg.SetUser(user)
	fmt.Println("New user created:", user)
	fmt.Printf("%+v\n", res)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAll(context.Background())
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	fmt.Println("All data erased successfully")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	res, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	if res == nil {
		return errors.New("no user available")
	}
	var output string
	for _, name := range res {
		if name == s.cfg.CurrentUserName {
			output += fmt.Sprintf("* %s (current)\n", name)
		} else {
			output += fmt.Sprintf("* %s\n", name)
		}
	}
	fmt.Print(output)
	return nil
}

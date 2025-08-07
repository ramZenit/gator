package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ramZenit/gator/internal/database"
)

func handlerListUsers(s *state, cmd command) error {
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

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ramZenit/gator/internal/database"
)

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

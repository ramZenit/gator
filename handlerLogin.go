package main

import (
	"context"
	"errors"
	"fmt"
)

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

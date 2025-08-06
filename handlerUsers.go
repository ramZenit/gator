package main

import (
	"context"
	"errors"
	"fmt"
)

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

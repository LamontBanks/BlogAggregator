package main

import (
	"context"
	"database/sql"
	"fmt"
)

func loginCommandInfo() commandInfo {
	return commandInfo{
		description: "Logs in a user",
		usage:       "login USERNAME",
		examples: []string{
			"login alice",
		},
	}
}

// Log in the user
// User must alrerady be registered
func handlerLogin(s *state, cmd command) error {
	// Args: username
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <username>", cmd.name)
	}
	username := cmd.args[0]

	// Check if the user is registered in the db
	// If nothing is returned, stop
	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%v not registered", username)
	}
	if err != nil {
		panic(err)
	}

	// Otherwise, "log in" the user by writing their name to the config file
	s.config.CurrentUserName = username
	if err := s.config.SetConfig(); err != nil {
		return err
	}
	fmt.Printf("Logged in as %v\n", username)

	return nil
}

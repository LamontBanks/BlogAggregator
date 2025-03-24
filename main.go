package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/config"
	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"

	// Leading underscore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

// -- Structs

// Application state to be passed to the commands:
// Config, database connection, etc.
type state struct {
	config *config.Config
	db     *database.Queries
}

// CLI command
type command struct {
	name string
	args []string
}

// Maps commands -> handler functions
type commands struct {
	cmds map[string]func(*state, command) error
}

// --- Main

func main() {
	// Register the CLI commands
	appCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)

	// Initialize info for the application state
	// Config
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Database connection
	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db) // Use the SQLC wrapper database instead of the SQL db directly

	// Set state
	appState := state{
		config: &cfg,
		db:     dbQueries,
	}

	// Read the CLI args to take action
	// os.Args includes the program name, then the command, and (possibly) args
	if len(os.Args) < 2 {
		log.Fatal("not enough args provided - need <command> <args>")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Run command
	cmdErr := appCommands.run(&appState, command{
		name: cmdName,
		args: cmdArgs,
	})
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

// -- CLI Command Handlers

// Log in the user
// User must alrerady be registered
// Usage:
//
//	$ go run . login <username>
//	$ go run . login alice
func handlerLogin(s *state, cmd command) error {
	// Get needed args
	if len(cmd.args) < 1 {
		return fmt.Errorf("username required")
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

	// Otherwise, log in the user by writing their name to the config file
	s.config.CurrentUserName = username
	if err := s.config.SetConfig(); err != nil {
		return err
	}
	fmt.Printf("Logged in as %v\n", username)

	return nil
}

// Register a user on the server, then pdates the config with the user.
// Usage:
//
//	$ go run . register <username>
//	$ go run . register alice
func handlerRegister(s *state, cmd command) error {
	// Get needed args
	if len(cmd.args) < 1 {
		return fmt.Errorf("name required")
	}
	username := cmd.args[0]

	// Insert user
	queryResult, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created user %v: %v\n", username, queryResult)

	// Update the config as well
	return handlerLogin(s, cmd)
}

// DEV/TESTING ONLY
// Deletes all users
func handlerReset(s *state, cmd command) error {
	return s.db.Reset(context.Background())
}

// -- Command functions

// Adds a new CLI command
// Command name is normalized to lowercase.
// Returns an errors if the command with the same name already exists
func (c *commands) register(name string, f func(*state, command) error) error {
	name = strings.ToLower(name)

	_, exists := c.cmds[name]
	if exists {
		return fmt.Errorf("command already exists: %v", name)
	}

	c.cmds[name] = f

	return nil
}

// Runs the function mapped to the named command
func (c *commands) run(s *state, cmd command) error {
	// Search the mapping for the assoicated handler function
	handlerFunc, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %v", cmd.name)
	}

	return handlerFunc(s, cmd)
}

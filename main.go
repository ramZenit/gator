package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/ramZenit/gator/internal/config"
	"github.com/ramZenit/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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
	cmdHandlers.register("agg", handlerAggregator)
	cmdHandlers.register("addfeed", handlerAddFeed)
	cmdHandlers.register("feeds", handlerFeeds)
	cmdHandlers.register("follow", handlerCreateFollow)
	cmdHandlers.register("following", handlerFollowsPerUser)

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

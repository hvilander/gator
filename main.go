package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/hvilander/gator/internal/config"
	"github.com/hvilander/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Config{}

	err := cfg.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("totally borked: %w", err))
		return
	}

	st := state{config: &cfg}
	cmds := commands{
		handlersByName: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	// db connection
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)
	st.db = dbQueries

	args := os.Args
	if len(args) < 2 {
		fmt.Println("you need to supply a command")
		os.Exit(1)
	}

	name := args[1]

	cmd := command{
		name: name,
		args: args[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Printf("run error: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)

}

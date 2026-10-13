package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/user14556454/project2/internal/database"

	_ "github.com/lib/pq"
	"github.com/user14556454/project2/internal/config"
)

func main() {
	existingConfig, err := config.Read()
	if err != nil {
		fmt.Println("Error in reading existing json:", err)
		os.Exit(1)
	}
	dbURL := existingConfig.DatabaseURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	sharedState := config.State{
		Db:           dbQueries,
		SharedConfig: &existingConfig,
	}

	cliCommands := config.Commands{
		CommandsMap: map[string]func(*config.State, config.Command) error{},
	}
	cliCommands.Register("login", config.HandlerLogin)
	cliCommands.Register("register", config.HandlerRegister)
	cliCommands.Register("reset", config.HandlerReset)
	cliCommands.Register("users", config.HandlerGetAllUsers)
	cliCommands.Register("agg", config.Handleragg)
	cliCommands.Register("addfeed", config.MiddlewareLoggedIn(config.Handleraddfeed))
	cliCommands.Register("feeds", config.Handlerfeeds)
	cliCommands.Register("follow", config.MiddlewareLoggedIn(config.Handlerfollow))
	cliCommands.Register("following", config.MiddlewareLoggedIn(config.Handlerfollowing))
	cliCommands.Register("unfollow", config.MiddlewareLoggedIn(config.Handlerunfollow))

	inputCommad := os.Args
	if len(inputCommad) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}
	commandinst := config.Command{
		Name:      inputCommad[1],
		Arguments: inputCommad[1:],
	}
	if err := cliCommands.Run(&sharedState, commandinst); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

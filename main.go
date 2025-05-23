package main

import (
	"database/sql"
	"fmt"
	commands "github.com/Cacutss/gator/internal/commands"
	config "github.com/Cacutss/gator/internal/config"
	database "github.com/Cacutss/gator/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

func parseArgs(args []string) []string {
	result := []string{}
	for _, v := range args {
		result = append(result, strings.Trim(v, "\""))
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Version 1.0.0")
		return
	}
	args := parseArgs(os.Args)
	state := config.State{}
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	state.Config = &conf
	db, err := sql.Open("postgres", state.Config.Dburl)
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries
	Commands := commands.GetCommands()
	command := commands.Command{
		Name: args[1],
		Args: args[1:],
	}
	if h, ok := Commands.Handler[args[1]]; ok {
		if err := h(&state, command); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unknown command")
	}
}

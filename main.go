package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Young-Pancake/goBlogGator/internal/config"
	"github.com/Young-Pancake/goBlogGator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	funcNames map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	
	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error reading database: %v", err)
	}
	dbQueries := database.New(db)
	st := state{
		db: dbQueries,
	}
	st.cfg = &cfg

	cmds := commands{funcNames: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("too few args")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatalf("error running config: %v", err)
	}
}

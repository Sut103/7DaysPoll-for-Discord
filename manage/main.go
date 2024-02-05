package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func createNewSession() (*discordgo.Session, error) {
	var input string
	fmt.Print("Enter Bot Token: ")
	fmt.Scan(&input)

	token := fmt.Sprintf("%s %s", "Bot", input)
	s, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	s.AddHandler(Ready)
	err = s.Open()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func run() error {

	if len(os.Args) < 2 {
		log.Fatalln("no command mode.")
	}

	switch os.Args[1] {
	case "register":
		return register()
	case "delete":
		return delete()
	default:
		log.Fatalln("unknown command mode.")
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

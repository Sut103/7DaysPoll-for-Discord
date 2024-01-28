package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Current Bot: %s#%s", s.State.User.Username, s.State.User.Discriminator)
}

func register() error {
	s, err := createNewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	// slash command
	command := discordgo.ApplicationCommand{
		Name:        "poll",
		Description: "Starting from today.",
	}

	// register
	ccmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &command)
	if err != nil {
		return err
	}

	fmt.Printf("Command registed: %+v\n", ccmd)
	return nil
}

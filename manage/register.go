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
	minLength := 5
	command := discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "poll",
		Description: "Starting 7DaysPoll from initial date (Today or Specific date).",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Please enter the title of the poll.",
				Type:        discordgo.ApplicationCommandOptionString,
			},
			{
				Name:        "start-date",
				Description: "If you have desired options, please specify the initial date. Example: 08/31",
				Type:        discordgo.ApplicationCommandOptionString,
				MaxLength:   5,
				MinLength:   &minLength,
			},
		},
	}

	// register
	ccmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", &command)
	if err != nil {
		return err
	}

	fmt.Printf("Command registed: %+v\n", ccmd)
	return nil
}

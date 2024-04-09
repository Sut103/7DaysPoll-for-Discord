package manage

import (
	"7DaysPoll/command"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Register(session *discordgo.Session) error {
	delete(session)
	
	log.Printf("Registering commands...\n")
	_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.GetPollCommand())
	if err != nil {
		return err
	}

	log.Printf("Command registration completed successfully.\n")
	return nil
}

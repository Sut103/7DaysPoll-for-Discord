package manage

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"7DaysPoll/internal/poll"
)

func Register(session *discordgo.Session) error {
	delete(session)
	log.Printf("Registering commands...\n")
	_, err := session.ApplicationCommandCreate(session.State.User.ID, "", poll.GetPollCommand())
	if err != nil {
		return err
	}
	log.Printf("Command registration completed successfully.\n")
	return nil
}

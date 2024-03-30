package manage

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// func Ready(session *discordgo.Session, r *discordgo.Ready) {
// 	log.Printf("Current Bot: %s#%s", session.State.User.Username, session.State.User.Discriminator)
// }

func Register(session *discordgo.Session) error {
	// slash command
	minLength := 5
	command := discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "poll",
		Description: "Starting 7DaysPoll from initial date (Today or Specify date).",
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

	// delete
	delete(session)

	// register
	log.Printf("Registering commands...\n")
	_, err := session.ApplicationCommandCreate(session.State.User.ID, "", &command)
	if err != nil {
		return err
	}

	log.Printf("Registering commands successfully completed.\n")
	return nil
}

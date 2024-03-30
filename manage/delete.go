package manage

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func delete(session *discordgo.Session) error {
	log.Printf("Deleting registed commands ...\n")
	commands, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return err
	}

	if len(commands) < 1 {
		log.Println("Could not find commands")
		return nil
	}

	for _, command := range commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "", command.ID)
		if err != nil {
			return err
		}
	}
	log.Printf("Deletion completed successfully. Deleted %d commands\n", len(commands))
	return nil
}

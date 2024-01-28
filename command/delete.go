package main

import (
	"fmt"
	"log"
)

func delete() error {
	s, err := createNewSession()
	if err != nil {
		return err
	}

	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return err
	}

	if len(commands) < 1 {
		log.Println("Could not find commands")
		return nil
	}

	for _, command := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", command.ID)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Successfully. Deleted %d commands\n", len(commands))
	return nil
}

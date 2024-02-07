package runtime

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"7DaysPoll-interactions/command"

	"github.com/bwmarrin/discordgo"
)

func poll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	body, err := command.Poll(i.Interaction)
	if err != nil {
		log.Println(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, body)
	if err != nil {
		log.Println(err)
	}
}

func RunBot(token string) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", "Bot", token))
	if err != nil {
		return err
	}

	s.AddHandler(poll)
	err = s.Open()
	if err != nil {
		return err
	}
	defer s.Close()

	log.Println("=====start=====")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	select {
	case <-signalChan:
		return nil
	}
}

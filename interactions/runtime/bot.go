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

func botHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := command.BotPoll(s, i)
	if err != nil {
		log.Println(err)
		return
	}
}

func RunBot(token string) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", "Bot", token))
	if err != nil {
		return err
	}

	s.AddHandler(botHandler)
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

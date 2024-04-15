package runtime

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"7DaysPoll/command"
	"7DaysPoll/manage"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	token string
}

func NewBot(token string) *Bot {
	return &Bot{
		token,
	}
}

func botHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := command.Poll(s, i.Interaction)
	if err != nil {
		log.Println(err)
		return
	}
}

func (b *Bot) Run() error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", "Bot", b.token))
	if err != nil {
		return err
	}

	s.AddHandler(botHandler)
	err = s.Open()
	if err != nil {
		return err
	}
	defer s.Close()
	
	manage.Register(s)

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

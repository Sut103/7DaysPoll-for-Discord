package runtime

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"7DaysPoll-interactions/command"
	"7DaysPoll-interactions/contextsync"

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

func messageReactionAddEventHandler(s *discordgo.Session, event *discordgo.MessageReactionAdd) {
	ctx := contextsync.NewAggregationContext()
	err := command.AggregatePoll(ctx, s, event.MessageReaction)
	if err != nil {
		log.Println(err)
		return
	}
}
func messageReactionRemoveEventHandler(s *discordgo.Session, event *discordgo.MessageReactionRemove) {
	ctx := contextsync.NewAggregationContext()
	err := command.AggregatePoll(ctx, s, event.MessageReaction)
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
	s.AddHandler(messageReactionAddEventHandler)
	s.AddHandler(messageReactionRemoveEventHandler)
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

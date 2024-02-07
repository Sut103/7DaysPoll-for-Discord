package runtime

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func RunBot(token string) error {
	s, err := discordgo.New(fmt.Sprintf("%s %s", "Bot", token))
	if err != nil {
		return err
	}

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

package command

import (
	"github.com/bwmarrin/discordgo"
)

type Session interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
}

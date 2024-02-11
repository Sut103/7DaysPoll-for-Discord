package command

import (
	"github.com/bwmarrin/discordgo"
)

type Session interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	InteractionResponse(interaction *discordgo.Interaction, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MessageReactionAdd(channelID, messageID, emojiID string, options ...discordgo.RequestOption) error
}

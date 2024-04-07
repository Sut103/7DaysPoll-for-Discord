package command

import (
	"github.com/bwmarrin/discordgo"
)

type Session interface {
	ChannelMessage(channelID, messageID string, options ...discordgo.RequestOption) (st *discordgo.Message, err error)
	ChannelMessageEditEmbeds(channelID, messageID string, embeds []*discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	InteractionResponse(interaction *discordgo.Interaction, options ...discordgo.RequestOption) (*discordgo.Message, error)
	MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string, options ...discordgo.RequestOption) ([]*discordgo.User, error)
	MessageReactionAdd(channelID, messageID, emojiID string, options ...discordgo.RequestOption) error
}

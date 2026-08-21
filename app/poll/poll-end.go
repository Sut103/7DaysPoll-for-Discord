package poll

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// GetPollEndMessageCommand registers a message context-menu command ("Apps"
// -> right-click a poll message) that lets a member end that poll early.
func GetPollEndMessageCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type: discordgo.MessageApplicationCommand,
		Name: "End Poll",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.Japanese: "投票を終了",
		},
	}
}

// EndPollMessageCommand expires the poll attached to the message the
// context-menu command was invoked on.
func EndPollMessageCommand(session *discordgo.Session, interaction *discordgo.Interaction) error {
	i18n := GetI18n(interaction.Locale)
	data := interaction.ApplicationCommandData()

	message, ok := data.Resolved.Messages[data.TargetID]
	if !ok {
		return respondEphemeral(session, interaction, i18n.PollNotFound)
	}
	return endPollMessage(session, interaction, message, i18n)
}

// GetPollEndSlashCommand registers the "/poll-end" slash command, an
// alternative to the "End Poll" message command for users who don't have
// (or don't want to use) the right-click context menu.
func GetPollEndSlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "poll-end",
		Description: "End a poll early. Only works on a poll message in this channel.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:         "message",
				Description:  "The poll message's link or ID (must be in this channel).",
				Type:         discordgo.ApplicationCommandOptionString,
				Required:     true,
				Autocomplete: true,
			},
		},
	}
}

// EndPollSlashCommand is the "/poll-end" handler. It resolves the "message"
// option to a message in the current channel and delegates to the same
// endPollMessage logic the "End Poll" context-menu command uses.
func EndPollSlashCommand(session *discordgo.Session, interaction *discordgo.Interaction) error {
	i18n := GetI18n(interaction.Locale)
	options := interaction.ApplicationCommandData().Options

	messageID, err := resolvePollMessageID(options[0].StringValue(), interaction.ChannelID)
	if err != nil {
		return respondEphemeral(session, interaction, i18n.PollEndInvalidMessage)
	}
	message, err := session.ChannelMessage(interaction.ChannelID, messageID)
	if err != nil {
		return respondEphemeral(session, interaction, i18n.PollNotFound)
	}
	return endPollMessage(session, interaction, message, i18n)
}

// pollEndAutocompleteScanLimit is how many of the channel's most recent
// messages are scanned for poll-end autocomplete suggestions, in a single
// ChannelMessages call (Discord's per-request maximum).
const pollEndAutocompleteScanLimit = 100

// PollEndAutocomplete answers the autocomplete interaction for /poll-end's
// "message" option. It scans the most recent messages in the invoking
// channel for not-yet-finalized poll messages, matches their question text
// against the user's partial input, and returns up to Discord's 25-choice
// limit — displaying the poll's title while the value sent on selection is
// the message ID (consumed unchanged by resolvePollMessageID).
func PollEndAutocomplete(session *discordgo.Session, interaction *discordgo.Interaction) error {
	data := interaction.ApplicationCommandData()
	query := strings.ToLower(strings.TrimSpace(data.Options[0].StringValue()))

	messages, err := session.ChannelMessages(interaction.ChannelID, pollEndAutocompleteScanLimit, "", "", "")
	if err != nil {
		log.Println("Failed to list channel messages for poll-end autocomplete:", err)
		return respondAutocomplete(session, interaction, nil)
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, message := range messages {
		if message.Poll == nil {
			continue
		}
		if message.Poll.Results != nil && message.Poll.Results.Finalized {
			continue
		}
		title := message.Poll.Question.Text
		if query != "" && !strings.Contains(strings.ToLower(title), query) {
			continue
		}
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  truncateRunes(title, 100),
			Value: message.ID,
		})
		if len(choices) == 25 {
			break
		}
	}
	return respondAutocomplete(session, interaction, choices)
}

func respondAutocomplete(session *discordgo.Session, interaction *discordgo.Interaction, choices []*discordgo.ApplicationCommandOptionChoice) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

var (
	messageIDPattern   = regexp.MustCompile(`^\d{17,20}$`)
	messageLinkPattern = regexp.MustCompile(`^https?://(?:canary\.|ptb\.)?(?:discord|discordapp)\.com/channels/(?:\d+|@me)/(\d+)/(\d+)$`)

	errInvalidMessageReference = errors.New("invalid poll message reference")
)

// resolvePollMessageID accepts either a bare Discord message ID or a full
// message link (https://discord.com/channels/{guild|@me}/{channel}/{message},
// including canary/ptb/discordapp.com variants) and returns the message ID.
// Links pointing at a channel other than currentChannelID are rejected,
// since canEndPoll's permission check below is only valid for the invoking
// channel (interaction.Member.Permissions is computed for that channel).
func resolvePollMessageID(input, currentChannelID string) (string, error) {
	input = strings.TrimSpace(input)
	if messageIDPattern.MatchString(input) {
		return input, nil
	}
	match := messageLinkPattern.FindStringSubmatch(input)
	if match == nil {
		return "", errInvalidMessageReference
	}
	channelID, messageID := match[1], match[2]
	if channelID != currentChannelID {
		return "", errInvalidMessageReference
	}
	return messageID, nil
}

// endPollMessage contains the guard checks and PollExpire call shared by the
// "End Poll" message command and the "/poll-end" slash command. Only the
// member who originally posted the poll, or a member with "Manage Messages"
// in the channel, may end it early — mirroring who Discord's own UI allows
// to end a native poll. (Discord's REST API only blocks ending a poll owned
// by a *different application*; per-member access control within our own
// bot's polls is entirely our responsibility.)
func endPollMessage(session *discordgo.Session, interaction *discordgo.Interaction, message *discordgo.Message, i18n I18n) error {
	if message.Poll == nil {
		return respondEphemeral(session, interaction, i18n.PollNotFound)
	}
	if message.Poll.Results != nil && message.Poll.Results.Finalized {
		return respondEphemeral(session, interaction, i18n.PollAlreadyEnded)
	}
	if !canEndPoll(interaction, message) {
		return respondEphemeral(session, interaction, i18n.PollEndNoPermission)
	}

	if _, err := session.PollExpire(message.ChannelID, message.ID); err != nil {
		log.Println("Failed to expire poll:", err)
		return respondEphemeral(session, interaction, i18n.PollEndFailed)
	}
	return respondEphemeral(session, interaction, i18n.PollEndSuccess)
}

// pollCreatorID returns the ID of the human who ran the slash command that
// produced this poll message. message.Author is the bot (it sent the
// interaction response), so the real invoker is recorded in
// message.InteractionMetadata.User instead (Message.Interaction is
// deprecated by Discord in favor of InteractionMetadata, so it is not used).
func pollCreatorID(message *discordgo.Message) string {
	if message.InteractionMetadata != nil && message.InteractionMetadata.User != nil {
		return message.InteractionMetadata.User.ID
	}
	return ""
}

func canEndPoll(interaction *discordgo.Interaction, message *discordgo.Message) bool {
	if interaction.Member == nil || interaction.Member.User == nil {
		return false
	}
	if id := pollCreatorID(message); id != "" && id == interaction.Member.User.ID {
		return true
	}
	return interaction.Member.Permissions&discordgo.PermissionManageMessages != 0
}

func respondEphemeral(session *discordgo.Session, interaction *discordgo.Interaction, content string) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

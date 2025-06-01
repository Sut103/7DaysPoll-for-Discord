package poll

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Choice struct {
	Emoji string
	Name  string
}

func getDays(day time.Time, numDays int) []time.Time {
	days := make([]time.Time, numDays)
	for i := range days {
		days[i] = day.AddDate(0, 0, i)
	}
	return days
}

func getEmojis() []string {
	return []string{
		"1⃣",
		"2⃣",
		"3⃣",
		"4⃣",
		"5⃣",
		"6⃣",
		"7⃣",
		"❌",
	}
}

func getChoices(locale discordgo.Locale, startDate time.Time, numDays int) []Choice {
	days := getDays(startDate, numDays)
	emojis := getEmojis()
	i18n := GetI18n(locale)

	choices := []Choice{}
	for i := 0; i < numDays; i++ {
		choices = append(choices, Choice{
			Emoji: emojis[i],
			Name:  fmt.Sprintf("%s (%s)", days[i].Format("01/02"), i18n.Weekdays[days[i].Weekday()]),
		})
	}
	absence := Choice{
		Emoji: emojis[7],
		Name:  i18n.Absence,
	}
	choices = append(choices, absence)
	return choices
}

func GetPollCommand() *discordgo.ApplicationCommand {
	minLength := 5
	minDays := 2
	maxDays := 7
	return &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "poll",
		Description: "Starting Poll from initial date with specified number of days (2-7).",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Please enter the title of the poll.",
				Type:        discordgo.ApplicationCommandOptionString,
			},
			{
				Name:        "start-date",
				Description: "If you have desired options, please specify the initial date. Example: 08/31",
				Type:        discordgo.ApplicationCommandOptionString,
				MaxLength:   5,
				MinLength:   &minLength,
			},
			{
				Name:        "days",
				Description: "Number of days for the poll (2-7). Default is 7.",
				Type:        discordgo.ApplicationCommandOptionInteger,
				MinValue:    FloatPtr(float64(minDays)),
				MaxValue:    float64(maxDays),
			},
		},
	}
}

func Poll(session *discordgo.Session, interaction *discordgo.Interaction) error {
	// get timezone
	timezone, err := GetTimeZone(string(interaction.Locale))
	if err != nil {
		log.Println(http.StatusInternalServerError, "timezone error", err)
		return err
	}
	// get options
	options := interaction.ApplicationCommandData().Options
	optMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range options {
		optMap[opt.Name] = opt
	}
	title := ""
	if t, ok := optMap["title"]; ok {
		title = t.StringValue()
	}
	// Get number of days (default: 7)
	numDays := 7
	if d, ok := optMap["days"]; ok {
		numDays = int(d.IntValue())
		if numDays < 2 {
			numDays = 2
		} else if numDays > 7 {
			numDays = 7
		}
	}
	// judgement start date
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timezone)
	if date, ok := optMap["start-date"]; ok {
		yearDate := fmt.Sprintf("%d/%s", now.Year(), date.StringValue())
		yd, err := time.Parse("2006/01/02", yearDate)
		if err == nil {
			if start.After(yd) {
				yd = yd.AddDate(1, 0, 0)
			}
			start = yd.In(timezone)
		}
	}
	// create response
	content := ""
	choices := getChoices(interaction.Locale, start, numDays)
	for _, choice := range choices {
		content += fmt.Sprintf("%s %s\n", choice.Emoji, choice.Name)
	}
	embed := discordgo.MessageEmbed{
		Title:       title,
		Description: "",
		Color:       0x780676,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "",
				Value:  content,
				Inline: true,
			},
			{
				Name:   "",
				Value:  "☑️ 0",
				Inline: true,
			},
		},
	}
	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	}
	err = session.InteractionRespond(interaction, &body)
	if err != nil {
		log.Println(err)
		return err
	}
	message, err := session.InteractionResponse(interaction)
	if err != nil {
		return err
	}
	for _, choice := range choices {
		err = session.MessageReactionAdd(interaction.ChannelID, message.ID, choice.Emoji)
		if err != nil {
			return err
		}
	}
	return nil
}

func AggregatePoll(ctx context.Context, session *discordgo.Session, reaction *discordgo.MessageReaction) error {
	me, err := session.User("@me")
	if err != nil {
		return err
	}
	if reaction.UserID == me.ID {
		return nil
	}
	emojis := getEmojis()
	isTargetEmoji := false
	for _, e := range emojis {
		if e == reaction.Emoji.Name {
			isTargetEmoji = true
			break
		}
	}
	if !isTargetEmoji {
		return nil
	}
	message, err := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return err
	}
	if !(len(message.Embeds) > 0 && len(message.Embeds[0].Fields) > 1) {
		return nil
	}
	go func() {
		embeds := message.Embeds
		embeds[0].Fields[1].Value = "☑️ ⌛" // It takes about 5 seconds for MessageReactions()
		session.ChannelMessageEditEmbeds(reaction.ChannelID, message.ID, embeds)
	}()
	uniqueVoter := map[string]struct{}{}
	time.Sleep(1 * time.Second)
	for _, e := range emojis {
		select {
		case <-ctx.Done():
			return nil
		default:
			func(emojiName string) {
				users, _ := session.MessageReactions(reaction.ChannelID, message.ID, emojiName, 100, "", "")
				for _, user := range users {
					uniqueVoter[user.ID] = struct{}{}
				}
			}(e)
		}
	}
	embeds := message.Embeds
	embeds[0].Fields[1].Value = fmt.Sprintf("☑️ %d", len(uniqueVoter)-1)
	session.ChannelMessageEditEmbeds(reaction.ChannelID, message.ID, embeds)
	return nil
}

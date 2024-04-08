package command

import (
	"7DaysPoll-interactions/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

func get7Days(day time.Time) [7]time.Time {
	days := [7]time.Time{}

	for i := range days {
		days[i] = day.AddDate(0, 0, i)
	}

	return days
}

func get7Emojis() []string {
	return []string{
		"1⃣",
		"2⃣",
		"3⃣",
		"4⃣",
		"5⃣",
		"6⃣",
		"7⃣",
	}
}

func Poll(session Session, interaction *discordgo.Interaction) error {
	// get timezone
	timezone, err := util.GetTimeZone(string(interaction.Locale))
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

	// prepare resources
	days := get7Days(start)
	emojis := get7Emojis()
	weekdays := util.GetWeekdays(interaction.Locale)

	// create response
	content := ""
	for i, day := range days {
		emoji := emojis[i]
		content += fmt.Sprintf("%s %s (%s)\n", emoji, day.Format("01/02"), weekdays[day.Weekday()])
	}

	embed := discordgo.MessageEmbed{
		Title:       title,
		Description: content,
		Color:       0x780676,
	}

	uniqueVoterCounter := discordgo.MessageEmbed{
		Title:       "",
		Description: "☑️ 0",
		Color:       0x780676,
	}

	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed, &uniqueVoterCounter},
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

	for _, reaction := range get7Emojis() {
		err = session.MessageReactionAdd(interaction.ChannelID, message.ID, reaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func AggregatePoll(ctx context.Context, session Session, reaction *discordgo.MessageReaction) error {
	me, err := session.User("@me")
	if err != nil {
		return err
	}
	if reaction.UserID == me.ID {
		return nil
	}

	message, err := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return err
	}

	go func() {
		embeds := append([]*discordgo.MessageEmbed{}, message.Embeds[0], &discordgo.MessageEmbed{
			Title:       "",
			Description: "☑️ ⌛", // It takes about 5 seconds for MessageReactions()
			Color:       0x780676,
		})
		session.ChannelMessageEditEmbeds(reaction.ChannelID, message.ID, embeds)
	}()

	uniqueVoter := map[string]struct{}{}

	time.Sleep(1 * time.Second)
	for _, r := range message.Reactions {
		select {
		case <-ctx.Done():
			return nil
		default:
			func(emojiName string) {
				// slow
				users, _ := session.MessageReactions(reaction.ChannelID, message.ID, emojiName, 100, "", "")
				for _, user := range users {
					uniqueVoter[user.ID] = struct{}{}
				}
			}(r.Emoji.Name)
		}
	}

	embeds := append([]*discordgo.MessageEmbed{}, message.Embeds[0], &discordgo.MessageEmbed{
		Title:       "",
		Description: fmt.Sprintf("☑️ %d", len(uniqueVoter)-1),
		Color:       0x780676,
	})
	session.ChannelMessageEditEmbeds(reaction.ChannelID, message.ID, embeds)

	return nil
}

package command

import (
	"7DaysPoll-interactions/util"
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

func Poll(interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	// get timezone
	timezone, err := util.GetTimeZone(string(interaction.Locale))
	if err != nil {
		log.Println(http.StatusInternalServerError, "timezone error", err)
		return nil, err
	}

	// get options
	options := interaction.ApplicationCommandData().Options
	optMap := map[string]*discordgo.ApplicationCommandInteractionDataOption{}
	for _, opt := range options {
		optMap[opt.Name] = opt
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
		Title:       "",
		Description: content,
		Color:       0x780676,
	}

	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	}
	return &body, nil
}

func BotPoll(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	body, err := Poll(i.Interaction)
	if err != nil {
		return err
	}

	err = s.InteractionRespond(i.Interaction, body)
	if err != nil {
		return err
	}

	message, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}

	for _, reaction := range get7Emojis() {
		err = s.MessageReactionAdd(i.Interaction.ChannelID, message.ID, reaction)
		if err != nil {
			return err
		}
	}
	return nil
}

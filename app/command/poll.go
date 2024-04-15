package command

import (
	"7DaysPoll/util"
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

func GetPollCommand() *discordgo.ApplicationCommand {
	minLength := 5
	return &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "poll",
		Description: "Starting 7DaysPoll from initial date (Today or Specific date).",
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
		},
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

	for _, reaction := range get7Emojis() {
		err = session.MessageReactionAdd(interaction.ChannelID, message.ID, reaction)
		if err != nil {
			return err
		}
	}
	return nil
}

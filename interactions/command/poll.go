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

	start := time.Now().Local().In(timezone)
	if date, ok := optMap["start-date"]; ok {
		d, err := time.Parse("20060102", date.StringValue())
		if err == nil {
			start = d.In(timezone)
		}
	}

	days := get7Days(start)
	emojis := get7Emojis()

	content := ""
	for i, day := range days {
		emoji := emojis[i]
		content += fmt.Sprintf("%s %d/%d/%d\n", emoji, day.Year(), day.Month(), day.Day())
	}

	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}
	return &body, nil
}

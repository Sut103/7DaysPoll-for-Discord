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
	timezone, err := util.GetTimeZone(string(interaction.Locale))
	if err != nil {
		log.Println(http.StatusInternalServerError, "timezone error", err)
		return nil, err
	}

	days := get7Days(time.Now().Local().In(timezone))
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

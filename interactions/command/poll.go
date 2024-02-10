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

func Poll(session Session, interaction *discordgo.Interaction) error {
	timezone, err := util.GetTimeZone(string(interaction.Locale))
	if err != nil {
		log.Println(http.StatusInternalServerError, "timezone error", err)
		return err
	}

	days := get7Days(time.Now().Local().In(timezone))
	emojis := get7Emojis()
	weekdays := util.GetWeekdays(interaction.Locale)

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

	err = session.InteractionRespond(interaction, &body)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

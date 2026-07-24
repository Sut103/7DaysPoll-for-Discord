package poll

import (
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

func getChoices(i18n I18n, startDate time.Time, numDays int) []Choice {
	days := getDays(startDate, numDays)
	emojis := getEmojis()

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

type pollOptions struct {
	Title   string
	Start   time.Time
	NumDays int
	OptMap  map[string]*discordgo.ApplicationCommandInteractionDataOption
}

func parsePollOptions(interaction *discordgo.Interaction, i18n I18n) (*pollOptions, error) {
	// get timezone
	timezone, err := GetTimeZone(string(interaction.Locale))
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
	title := ""
	if t, ok := optMap["title"]; ok {
		title = t.StringValue()
	}
	if title == "" {
		title = i18n.DefaultTitle
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
	return &pollOptions{
		Title:   title,
		Start:   start,
		NumDays: numDays,
		OptMap:  optMap,
	}, nil
}

func buildMessageURL(guildID, channelID, messageID string) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID)
}

func buildEventURL(guildID, eventID string) string {
	return fmt.Sprintf("https://discord.com/events/%s/%s", guildID, eventID)
}

const discordEventNameMaxLength = 100

func createScheduledEvent(session *discordgo.Session, guildID string, i18n I18n, start time.Time, numDays int, title string, messageURL string) (*discordgo.GuildScheduledEvent, error) {
	eventTitle := truncateRunes(i18n.VotingPeriod+title, discordEventNameMaxLength)

	days := getDays(start, numDays)
	finalDay := days[len(days)-1]

	// Use the final day of the voting period as the event start time, since once
	// the scheduled start time passes the event begins automatically and its
	// start time can no longer be updated after the date is decided.
	startTime := time.Date(finalDay.Year(), finalDay.Month(), finalDay.Day(), 0, 0, 0, 0, start.Location())
	now := time.Now()
	// Discord API requires scheduled start time to be in the future
	if startTime.Before(now) {
		startTime = now.Add(1 * time.Minute)
	}
	endTime := time.Date(finalDay.Year(), finalDay.Month(), finalDay.Day(), 23, 59, 59, 0, start.Location())

	eventParams := &discordgo.GuildScheduledEventParams{
		Name:               eventTitle,
		Description:        fmt.Sprintf("%s: %s", i18n.PollMessage, messageURL),
		ScheduledStartTime: &startTime,
		ScheduledEndTime:   &endTime,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: messageURL,
		},
	}

	return session.GuildScheduledEventCreate(guildID, eventParams)
}

package util

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetTimeZone(lang string) (*time.Location, error) {
	timezone := map[string]string{
		"Japanese": "Asia/Tokyo",
	}

	tz, ok := timezone[lang]
	if !ok {
		return time.Local, nil
	}

	return time.LoadLocation(tz)
}

func GetWeekdays(lang discordgo.Locale) []string {
	localeWeekdays := map[discordgo.Locale][]string{
		discordgo.EnglishUS: {"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		discordgo.Japanese:  {"日", "月", "火", "水", "木", "金", "土"},
	}

	weekdays, ok := localeWeekdays[lang]
	if !ok {
		return localeWeekdays[discordgo.EnglishUS]
	}

	return weekdays
}

func GetAbsence(lang discordgo.Locale) string {
	absence := map[discordgo.Locale]string{
		discordgo.EnglishUS: "Absence",
		discordgo.Japanese:  "欠席",
	}

	name, ok := absence[lang]
	if !ok {
		return absence[discordgo.EnglishUS]
	}

	return name
}

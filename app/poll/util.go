package poll

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

func getWeekdays(lang discordgo.Locale) []string {
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

func getAbsence(lang discordgo.Locale) string {
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

type I18n struct {
	Weekdays []string
	Absence  string
}

func GetI18n(lang discordgo.Locale) I18n {
	return I18n{
		Weekdays: getWeekdays(lang),
		Absence:  getAbsence(lang),
	}
}

func FloatPtr(v float64) *float64 {
	return &v
}

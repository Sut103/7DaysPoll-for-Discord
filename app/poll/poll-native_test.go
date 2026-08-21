package poll

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestGetNativePollCommand(t *testing.T) {
	cmd := GetNativePollCommand()

	if cmd.Name != "poll" {
		t.Errorf("Name = %q, want %q", cmd.Name, "poll")
	}
	if cmd.Type != discordgo.ChatApplicationCommand {
		t.Errorf("Type = %v, want %v", cmd.Type, discordgo.ChatApplicationCommand)
	}

	optionsByName := map[string]*discordgo.ApplicationCommandOption{}
	for _, opt := range cmd.Options {
		optionsByName[opt.Name] = opt
	}

	titleOpt, ok := optionsByName["title"]
	if !ok {
		t.Fatal("expected a \"title\" option")
	}
	if titleOpt.Type != discordgo.ApplicationCommandOptionString {
		t.Errorf("title option Type = %v, want string", titleOpt.Type)
	}
	if titleOpt.MaxLength != pollQuestionMaxLength {
		t.Errorf("title option MaxLength = %d, want %d (Discord's poll question limit)", titleOpt.MaxLength, pollQuestionMaxLength)
	}

	daysOpt, ok := optionsByName["days"]
	if !ok {
		t.Fatal("expected a \"days\" option")
	}
	if daysOpt.MinValue == nil || *daysOpt.MinValue != 2 {
		t.Errorf("days option MinValue = %v, want 2", daysOpt.MinValue)
	}
	if daysOpt.MaxValue != 7 {
		t.Errorf("days option MaxValue = %v, want 7", daysOpt.MaxValue)
	}

	durationOpt, ok := optionsByName["duration"]
	if !ok {
		t.Fatal("expected a \"duration\" option")
	}
	if durationOpt.MinValue == nil || *durationOpt.MinValue != minDurationDays {
		t.Errorf("duration option MinValue = %v, want %d", durationOpt.MinValue, minDurationDays)
	}
	if durationOpt.MaxValue != maxDurationDays {
		t.Errorf("duration option MaxValue = %v, want %d", durationOpt.MaxValue, maxDurationDays)
	}

	startDateOpt, ok := optionsByName["start-date"]
	if !ok {
		t.Fatal("expected a \"start-date\" option")
	}
	if startDateOpt.MinLength == nil || *startDateOpt.MinLength != 5 {
		t.Errorf("start-date option MinLength = %v, want 5", startDateOpt.MinLength)
	}
	if startDateOpt.MaxLength != 5 {
		t.Errorf("start-date option MaxLength = %d, want 5", startDateOpt.MaxLength)
	}
}

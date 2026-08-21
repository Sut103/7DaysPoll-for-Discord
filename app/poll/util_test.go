package poll

import (
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestGetTimeZone(t *testing.T) {
	t.Run("japanese locale maps to Asia/Tokyo", func(t *testing.T) {
		got, err := GetTimeZone(discordgo.Japanese)
		if err != nil {
			t.Fatalf("GetTimeZone(Japanese) returned error: %v", err)
		}
		if got.String() != "Asia/Tokyo" {
			t.Errorf("GetTimeZone(Japanese) = %v, want Asia/Tokyo", got.String())
		}
	})

	t.Run("unknown locale falls back to the process's time.Local", func(t *testing.T) {
		// Compared by identity, not by name: time.Local and time.UTC are
		// always distinct *Location values, so this catches the fallback
		// ever changing to a different fixed zone even in an environment
		// (e.g. CI, which is often TZ=UTC) where the two would otherwise
		// coincidentally have the same name.
		got, err := GetTimeZone(discordgo.EnglishUS)
		if err != nil {
			t.Fatalf("GetTimeZone(EnglishUS) returned error: %v", err)
		}
		if got != time.Local {
			t.Errorf("GetTimeZone(EnglishUS) = %v, want time.Local", got)
		}
	})
}

func TestGetWeekdays(t *testing.T) {
	tests := []struct {
		name   string
		locale discordgo.Locale
		want   []string
	}{
		{"english", discordgo.EnglishUS, []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}},
		{"japanese", discordgo.Japanese, []string{"日", "月", "火", "水", "木", "金", "土"}},
		{"unknown locale falls back to english", discordgo.German, []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getWeekdays(tt.locale)
			if len(got) != len(tt.want) {
				t.Fatalf("getWeekdays(%v) = %v, want %v", tt.locale, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getWeekdays(%v)[%d] = %q, want %q", tt.locale, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetAbsence(t *testing.T) {
	tests := []struct {
		name   string
		locale discordgo.Locale
		want   string
	}{
		{"english", discordgo.EnglishUS, "Absence"},
		{"japanese", discordgo.Japanese, "欠席"},
		{"unknown locale falls back to english", discordgo.German, "Absence"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAbsence(tt.locale); got != tt.want {
				t.Errorf("getAbsence(%v) = %q, want %q", tt.locale, got, tt.want)
			}
		})
	}
}

func TestGetTitle(t *testing.T) {
	tests := []struct {
		name   string
		locale discordgo.Locale
		want   string
	}{
		{"english", discordgo.EnglishUS, "Poll"},
		{"japanese", discordgo.Japanese, "投票"},
		{"unknown locale falls back to english", discordgo.German, "Poll"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTitle(tt.locale); got != tt.want {
				t.Errorf("getTitle(%v) = %q, want %q", tt.locale, got, tt.want)
			}
		})
	}
}

func TestGetVotingPeriod(t *testing.T) {
	tests := []struct {
		name   string
		locale discordgo.Locale
		want   string
	}{
		{"english", discordgo.EnglishUS, "(🗳️Voting)"},
		{"japanese", discordgo.Japanese, "(🗳️投票期間中)"},
		{"unknown locale falls back to english", discordgo.German, "(🗳️Voting)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVotingPeriod(tt.locale); got != tt.want {
				t.Errorf("getVotingPeriod(%v) = %q, want %q", tt.locale, got, tt.want)
			}
		})
	}
}

func TestGetPollMessage(t *testing.T) {
	tests := []struct {
		name   string
		locale discordgo.Locale
		want   string
	}{
		{"english", discordgo.EnglishUS, "Poll message"},
		{"japanese", discordgo.Japanese, "投票メッセージ"},
		{"unknown locale falls back to english", discordgo.German, "Poll message"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPollMessage(tt.locale); got != tt.want {
				t.Errorf("getPollMessage(%v) = %q, want %q", tt.locale, got, tt.want)
			}
		})
	}
}

func TestGetI18n(t *testing.T) {
	got := GetI18n(discordgo.Japanese)
	want := I18n{
		Weekdays:     []string{"日", "月", "火", "水", "木", "金", "土"},
		Absence:      "欠席",
		DefaultTitle: "投票",
		VotingPeriod: "(🗳️投票期間中)",
		PollMessage:  "投票メッセージ",
	}
	if got.Absence != want.Absence || got.DefaultTitle != want.DefaultTitle ||
		got.VotingPeriod != want.VotingPeriod || got.PollMessage != want.PollMessage {
		t.Fatalf("GetI18n(Japanese) = %+v, want %+v", got, want)
	}
	if strings.Join(got.Weekdays, ",") != strings.Join(want.Weekdays, ",") {
		t.Errorf("GetI18n(Japanese).Weekdays = %v, want %v", got.Weekdays, want.Weekdays)
	}
}

func TestFloatPtr(t *testing.T) {
	got := FloatPtr(3.5)
	if got == nil {
		t.Fatal("FloatPtr(3.5) returned nil")
	}
	if *got != 3.5 {
		t.Errorf("*FloatPtr(3.5) = %v, want 3.5", *got)
	}
}

func TestTruncateRunes(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"shorter than max is untouched", "hello", 10, "hello"},
		{"exactly at max is untouched", "hello", 5, "hello"},
		{"ascii longer than max is truncated", "hello world", 5, "hello"},
		// truncateRunes must count runes, not bytes: each of these Japanese
		// characters is 3 bytes in UTF-8, so a naive byte-slice truncation
		// would split a rune in half and corrupt the string.
		{"multi-byte string is truncated by rune count, not bytes", "投票期間中です", 3, "投票期"},
		{"multi-byte string shorter than max is untouched", "投票", 5, "投票"},
		{"maxLen zero yields empty string", "hello", 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateRunes(tt.s, tt.maxLen); got != tt.want {
				t.Errorf("truncateRunes(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
			}
		})
	}
}

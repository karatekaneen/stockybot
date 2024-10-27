package bot

import (
	"testing"
	"time"

	"github.com/karatekaneen/stockybot"
	"github.com/matryer/is"
)

func Test_daysSinceLast(t *testing.T) {
	tests := []struct {
		name    string
		signals []stockybot.Signal
		now     time.Time
		want    int
	}{
		{
			name:    "No signals",
			want:    0,
			signals: nil,
			now:     time.Now(),
		},
		{
			name: "Simple",
			want: 31,
			signals: []stockybot.Signal{
				{Date: time.Date(2023, 3, 24, 0, 24, 23, 0, time.Local)},
			},
			now: time.Date(2023, 4, 24, 23, 24, 23, 0, time.Local),
		},
		{
			name: "Multiple",
			want: 23,
			signals: []stockybot.Signal{
				{Date: time.Date(2023, 3, 2, 0, 24, 23, 0, time.Local)},
				{Date: time.Date(2021, 3, 12, 0, 24, 23, 0, time.Local)},
				{Date: time.Date(2022, 3, 4, 0, 24, 23, 0, time.Local)},
				{Date: time.Date(2023, 3, 24, 0, 24, 23, 0, time.Local)},
				{Date: time.Date(2023, 4, 1, 11, 23, 34, 0, time.Local)},
			},
			now: time.Date(2023, 4, 24, 23, 24, 23, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(daysSinceLast(tt.signals, tt.now), tt.want)
		})
	}
}

func Test_numberOfExitsSince(t *testing.T) {
	tests := []struct {
		limit   time.Time
		name    string
		signals []stockybot.Signal
		want    int
	}{
		{
			name: "Multiple",
			want: 2,
			signals: []stockybot.Signal{
				{
					Date: time.Date(2020, 3, 2, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2021, 3, 12, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2022, 8, 4, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2023, 3, 24, 0, 24, 23, 0, time.Local),
					Type: "enter", // This one is ignored
				},
				{
					Date: time.Date(2023, 4, 1, 11, 23, 34, 0, time.Local),
					Type: "exit",
				},
			},
			limit: time.Date(2023, 7, 2, 23, 24, 23, 0, time.Local).Add(-(time.Hour * 24 * 365)),
		},
		{
			name: "Multiple",
			want: 0,
			signals: []stockybot.Signal{
				{
					Date: time.Date(2020, 3, 2, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2021, 3, 12, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2022, 8, 4, 0, 24, 23, 0, time.Local),
					Type: "exit",
				},
				{
					Date: time.Date(2023, 3, 24, 0, 24, 23, 0, time.Local),
					Type: "enter", // This one is ignored
				},
				{
					Date: time.Date(2023, 4, 1, 11, 23, 34, 0, time.Local),
					Type: "exit",
				},
			},
			limit: time.Date(2023, 7, 2, 23, 24, 23, 0, time.Local).Add(-(time.Hour * 24)),
		},
		{
			name:    "No Signals",
			want:    0,
			signals: nil,
			limit:   time.Date(2023, 7, 2, 23, 24, 23, 0, time.Local).Add(-(time.Hour * 24)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(numberOfExitsSince(tt.signals, tt.limit), tt.want)
		})
	}
}

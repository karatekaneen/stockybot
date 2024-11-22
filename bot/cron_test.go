package bot

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/matryer/is"

	"github.com/karatekaneen/stockybot"
)

func Test_printSummary(t *testing.T) {
	is := is.New(t)

	summary := dailySummary{
		Sells: []watchSignal{
			{
				prediction: prediction{
					Signal: stockybot.Signal{
						Stock: stockybot.Security{
							List:    "Ankeborg Large Cap",
							Name:    "Ankeborgs Kod & Keb AB",
							Type:    "stock",
							Country: "Sweden",
							ID:      123,
						},
						Price: sql.NullFloat64{Float64: 1.2345, Valid: true},
					},
					score: 0,
				},
				Watchers: []string{"the_bear"},
			},
			{
				prediction: prediction{
					Signal: stockybot.Signal{
						Stock: stockybot.Security{
							List:    "Ankeborg Large Cap",
							Name:    "The company",
							Type:    "stock",
							Country: "Sweden",
							ID:      123,
						},
						Price: sql.NullFloat64{Float64: 1.2345, Valid: true},
					},
					score: 0,
				},
				Watchers: []string{},
			},
		},

		Buys: []watchSignal{
			{
				prediction: prediction{
					Signal: stockybot.Signal{
						Stock: stockybot.Security{
							List:    "Ankeborg Large Cap",
							Name:    "Ankeborgs Bank",
							Type:    "stock",
							Country: "Sweden",
							ID:      123,
						},
						Price: sql.NullFloat64{Float64: 1.2345, Valid: true},
					},
					score: 0.666666666667,
				},
				Watchers: []string{"kalleanka", "joakim_von_ca$h"},
			},
			{
				prediction: prediction{
					Signal: stockybot.Signal{
						Stock: stockybot.Security{
							List: "Ankeborg Large Cap",
							Name: "Ankeborgs Gold AB",
							ID:   123,
						},
						Price: sql.NullFloat64{Float64: 1.2345, Valid: true},
					},
					score: 0.50,
				},
				Watchers: []string{},
			},
			{
				prediction: prediction{
					Signal: stockybot.Signal{
						Stock: stockybot.Security{
							List:    "Ankeborg Large Cap",
							Name:    "Dunder Mifflin",
							Type:    "stock",
							Country: "Sweden",
							ID:      123,
						},
						Price: sql.NullFloat64{Float64: 1.2345, Valid: true},
					},
					score: 0.1,
				},
				Watchers: []string{"joakim_von_ca$h"},
			},
		},
	}

	got, err := printSummary(reportTemplate, summary)

	is.NoErr(err)
	fmt.Println(got)
	is.Equal(len(got), 261)
}

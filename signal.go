package stockybot

import (
	"database/sql"
	"time"
)

type Signal struct {
	Action      string          `firestore:"action,omitempty"`      // buy or sell
	Date        time.Time       `firestore:"date,omitempty"`        // The date for when the signal was executed. Or when it was generated for pending signals
	Price       sql.NullFloat64 `firestore:"price,omitempty"`       // The price where the signal was executed
	Status      string          `firestore:"status,omitempty"`      // pending or executed
	TriggerDate time.Time       `firestore:"triggerDate,omitempty"` // The date for when the signal was executed
	Type        string          `firestore:"type,omitempty"`        // Entry or exit. Used in combination with Action to define what kind of signal to be able to use short positions in the future.
	Stock       Security        `firestore:"stock,omitempty"`       // Only partial stock data loaded here
	Strategy    string          `firestore:"strategy,omitempty"`    // What strategy generated the signal
}

type Security struct {
	ID       int    `firestore:"id,omitempty"`       // ID of the security. Corresponds to the ID on Avanza
	List     string `firestore:"list,omitempty"`     // What list the stock is listed in
	Name     string `firestore:"name,omitempty"`     // Name of the stock
	LinkName string `firestore:"linkName,omitempty"` // URL safe name
	Type     string `firestore:"type,omitempty"`     // index or stock. Rates will be added some time in the future as well
	Country  string `firestore:"country,omitempty"`  // Country code. "SE" for Sweden
}

type PriceDocument struct {
	LastPricePoint time.Time    `firestore:"lastPricePoint,omitempty"`
	UpdatedAt      time.Time    `firestore:"updatedAt,omitempty"`
	PriceData      []PricePoint `firestore:"priceData,omitempty"`
	Type           string       `firestore:"type,omitempty"` // "stock" or "index"
}

type PricePoint struct {
	Close  float64 `firestore:"close,omitempty"`
	Open   float64 `firestore:"open,omitempty"`
	High   float64 `firestore:"high,omitempty"`
	Low    float64 `firestore:"low,omitempty"`
	Owners int     `firestore:"owners,omitempty"` // Number of owners on avanza
	Volume int     `firestore:"volume,omitempty"` // Number of shares traded
	Date   string  `firestore:"date,omitempty"`
}

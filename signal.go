package stockybot

import (
	"database/sql"
	"fmt"
	"time"
)

// Describes the current data that the strategy keeps
// about a particular stock
type StrategyState struct {
	LastSignal   *Signal `firestore:"lastSignal,omitempty"`
	Bias         string  `firestore:"bias,omitempty"`
	Regime       string  `firestore:"regime,omitempty"`
	HighPrice    float64 `firestore:"highPrice,omitempty"`
	LowPrice     float64 `firestore:"lowPrice,omitempty"`
	TriggerPrice float64 `firestore:"triggerPrice,omitempty"`
}

type Signal struct {
	Date        time.Time       `firestore:"date,omitempty"`
	TriggerDate time.Time       `firestore:"triggerDate,omitempty"`
	Action      string          `firestore:"action,omitempty"`
	Status      string          `firestore:"status,omitempty"`
	Type        string          `firestore:"type,omitempty"`
	Strategy    string          `firestore:"strategy,omitempty"`
	Stock       Security        `firestore:"stock,omitempty"`
	Price       sql.NullFloat64 `firestore:"price,omitempty"`
}

type Security struct {
	List     string `firestore:"list,omitempty"`
	Name     string `firestore:"name,omitempty"`
	LinkName string `firestore:"linkName,omitempty"`
	Type     string `firestore:"type,omitempty"`
	Country  string `firestore:"country,omitempty"`
	ID       int64  `firestore:"id,omitempty"`
}

type PriceDocument struct {
	UpdatedAt      time.Time    `firestore:"updatedAt,omitempty"`
	LastPricePoint string       `firestore:"lastPricePoint,omitempty"`
	Type           string       `firestore:"type,omitempty"`
	PriceData      []PricePoint `firestore:"priceData,omitempty"`
}

type PricePoint struct {
	Date   string  `json:"date,omitempty"   firestore:"date,omitempty"`
	Close  float64 `json:"close,omitempty"  firestore:"close,omitempty"`
	Open   float64 `json:"open,omitempty"   firestore:"open,omitempty"`
	High   float64 `json:"high,omitempty"   firestore:"high,omitempty"`
	Low    float64 `json:"low,omitempty"    firestore:"low,omitempty"`
	Volume int64   `json:"volume,omitempty" firestore:"volume,omitempty"`
}

// You will have a bad time if you pass a negative number. Don't do that.
func LastN[T any](arr []T, n int) ([]T, error) {
	if n > len(arr) {
		return arr, fmt.Errorf("list not long enough")
	}

	return arr[len(arr)-n:], nil
}

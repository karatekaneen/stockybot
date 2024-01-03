package firestore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/karatekaneen/stockybot"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type rawSignal struct {
	Action      string             `firestore:"action,omitempty"`      // buy or sell
	Date        string             `firestore:"date,omitempty"`        // The date for when the signal was executed. Or when it was generated for pending signals
	Price       sql.NullFloat64    `firestore:"price,omitempty"`       // The price where the signal was executed
	Status      string             `firestore:"status,omitempty"`      // pending or executed
	TriggerDate string             `firestore:"triggerDate,omitempty"` // The date for when the signal was executed
	Type        string             `firestore:"type,omitempty"`        // Entry or exit. Used in combination with Action to define what kind of signal to be able to use short positions in the future.
	Stock       stockybot.Security `firestore:"stock,omitempty"`       // Only partial stock data loaded here
	Strategy    string             `firestore:"strategy,omitempty"`    // What strategy generated the signal
}

type rawSignal2 struct {
	Action      string             `firestore:"action,omitempty"`      // buy or sell
	Date        string             `firestore:"date,omitempty"`        // The date for when the signal was executed. Or when it was generated for pending signals
	Price       float64            `firestore:"price,omitempty"`       // The price where the signal was executed
	Status      string             `firestore:"status,omitempty"`      // pending or executed
	TriggerDate string             `firestore:"triggerDate,omitempty"` // The date for when the signal was executed
	Type        string             `firestore:"type,omitempty"`        // Entry or exit. Used in combination with Action to define what kind of signal to be able to use short positions in the future.
	Stock       stockybot.Security `firestore:"stock,omitempty"`       // Only partial stock data loaded here
	Strategy    string             `firestore:"strategy,omitempty"`    // What strategy generated the signal
}

type signalDocument struct {
	Signals []rawSignal2 `firestore:"signals,omitempty"`
}

func (r rawSignal) toSignal() (*stockybot.Signal, error) {
	date, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		return nil, errors.Wrap(err, "Parse: ")
	}
	triggerDate, err := time.Parse(time.RFC3339, r.TriggerDate)
	if err != nil {
		return nil, errors.Wrap(err, "Parse: ")
	}

	return &stockybot.Signal{
		Action:      r.Action,
		Date:        date,
		Price:       r.Price,
		Status:      r.Status,
		TriggerDate: triggerDate,
		Type:        r.Type,
		Stock:       r.Stock,
		Strategy:    "flipper",
	}, nil
}

func (r rawSignal2) toSignal() (*stockybot.Signal, error) {
	date, err := time.Parse(time.RFC3339, r.Date)
	if err != nil {
		return nil, errors.Wrap(err, "Parse: ")
	}
	triggerDate, err := time.Parse(time.RFC3339, r.TriggerDate)
	if err != nil {
		return nil, errors.Wrap(err, "Parse: ")
	}

	return &stockybot.Signal{
		Action:      r.Action,
		Date:        date,
		Price:       sql.NullFloat64{Float64: r.Price, Valid: true},
		Status:      r.Status,
		TriggerDate: triggerDate,
		Type:        r.Type,
		Stock:       r.Stock,
		Strategy:    "flipper",
	}, nil
}

func (f *FireDB) PendingSignals(ctx context.Context) ([]stockybot.Signal, error) {
	signals := []stockybot.Signal{}

	signalIterator := f.client.Collection("pending-signals").Documents(ctx)

	for {
		doc, err := signalIterator.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		var raw rawSignal

		if err := doc.DataTo(&raw); err != nil {
			return nil, errors.Wrap(err, "conversion: ")
		}

		sig, err := raw.toSignal()
		if err != nil {
			return nil, err
		}

		signals = append(signals, *sig)
	}

	return signals, nil
}

func (f *FireDB) Signals(ctx context.Context, stockId int64) ([]stockybot.Signal, error) {
	signals := []stockybot.Signal{}

	doc, err := f.client.Collection("signals").Doc(fmt.Sprint(stockId)).Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching signal document")
	} else if !doc.Exists() {
		return nil, stockybot.ErrNotFound
	}

	var sigDoc signalDocument
	if err := doc.DataTo(&sigDoc); err != nil {
		return nil, errors.Wrap(err, "conversion: ")
	}

	for _, raw := range sigDoc.Signals {
		sig, err := raw.toSignal()
		if err != nil {
			return nil, err
		}

		signals = append(signals, *sig)
	}

	return signals, nil
}

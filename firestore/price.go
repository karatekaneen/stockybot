package firestore

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/karatekaneen/stockybot"
)

var ErrNotExist = errors.New("document does not exist")

func (f *FireDB) PriceData(ctx context.Context, id int64) (*stockybot.PriceDocument, error) {
	doc, err := f.client.Collection("prices").Doc(strconv.FormatInt(id, 10)).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch price document: %w", err)
	} else if !doc.Exists() {
		return nil, ErrNotExist
	}

	var raw stockybot.PriceDocument
	if err := doc.DataTo(&raw); err != nil {
		return nil, fmt.Errorf("convert price document: %w", err)
	}

	return &raw, nil
}

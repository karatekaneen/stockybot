package firestore

import (
	"context"
	"fmt"

	"github.com/karatekaneen/stockybot"
	"github.com/pkg/errors"
)

func (f *FireDB) PriceData(ctx context.Context, id int64) (*stockybot.PriceDocument, error) {
	doc, err := f.client.Collection("prices").Doc(fmt.Sprint(id)).Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Get document:")
	} else if !doc.Exists() {
		return nil, errors.New("Document does not exist")
	}

	var raw stockybot.PriceDocument
	if err := doc.DataTo(&raw); err != nil {
		return nil, errors.Wrap(err, "conversion: ")
	}

	return &raw, nil
}

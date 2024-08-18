package firestore

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/iterator"

	"github.com/karatekaneen/stockybot"
)

func (f *FireDB) Securities(ctx context.Context) ([]stockybot.Security, error) {
	securities := []stockybot.Security{}

	iter := f.client.Collection("securities").Documents(ctx)

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("fetch securities: %w", err)
		}

		var sec stockybot.Security
		if err := doc.DataTo(&sec); err != nil {
			return nil, fmt.Errorf("convert document %q to Security: %w", doc.Ref.Path, err)
		}

		securities = append(securities, sec)
	}

	return securities, nil
}

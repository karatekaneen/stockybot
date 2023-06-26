package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

type FireDB struct {
	client *firestore.Client
}

// createFirestoreClient creates a firestore client
func New(ctx context.Context, projectId string) (*FireDB, error) {
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, errors.Wrap(err, "firestore init:")
	}

	f := &FireDB{client}

	return f, nil
}

func (f *FireDB) Close() error {
	return f.client.Close()
}

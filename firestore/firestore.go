package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/karatekaneen/stockybot"
	"github.com/pkg/errors"
)

type FireDB struct {
	client *firestore.Client
}

type Config struct {
	ProjectID string `help:"GCP project that hosts the database" required:"" env:"GCP_PROJECT_ID"`
}

// Security implements bot.dataRepository.
func (*FireDB) Security(ctx context.Context, id int64) (*stockybot.Security, error) {
	panic("unimplemented")
}

// createFirestoreClient creates a firestore client
func New(ctx context.Context, cfg Config) (*FireDB, error) {
	client, err := firestore.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "firestore init:")
	}

	f := &FireDB{client}

	return f, nil
}

func (f *FireDB) Close() error {
	return f.client.Close()
}

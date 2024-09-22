package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/ent"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

type SecurityProvider interface {
	Securities(ctx context.Context) ([]stockybot.Security, error)
}

type DB struct {
	Client *ent.Client
	log    *zap.SugaredLogger
}

type Config struct {
	Location string `help:"DB file location" env:"DB_LOCATION"`
}

func (c Config) DSN() string {
	if c.Location == "" {
		return "file:ent?mode=memory&cache=shared&_fk=1"
	}

	return c.Location + "?cache=shared&_fk=1"
}

func New(ctx context.Context, cfg Config, log *zap.SugaredLogger) (*DB, error) {
	client, err := ent.Open("sqlite3", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &DB{Client: client, log: log}, nil
}

func (db *DB) Close() error {
	return db.Client.Close() //nolint:wrapcheck
}

func (db *DB) ImportPeriodically(
	ctx context.Context,
	provider SecurityProvider,
	cadence time.Duration,
) error {
	// Do a first import before starting timer
	if err := db.ImportSecurities(ctx, provider); err != nil {
		return fmt.Errorf("first import failed: %w", err)
	}

	if cadence == 0 {
		return nil
	}

	ticker := time.NewTicker(cadence)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err() //nolint:wrapcheck
		case <-ticker.C:
			db.log.Info("Starting security import")

			if err := db.ImportSecurities(ctx, provider); err != nil {
				return fmt.Errorf("import failed: %w", err)
			}
		}
	}
}

func (db *DB) ImportSecurities(ctx context.Context, provider SecurityProvider) error {
	secs, err := provider.Securities(ctx)
	if err != nil {
		return fmt.Errorf("fetch securities from provider: %w", err)
	}

	db.log.Infof("Found %d stocks to import", len(secs))

	created := 0

	for _, sec := range secs {
		err := db.Client.Security.Create().
			SetID(sec.ID).
			SetName(sec.Name).
			SetList(sec.List).
			SetLinkName(sec.LinkName).
			SetType(security.Type(sec.Type)).
			SetCountry(sec.Country).
			Exec(ctx)
		if ent.IsConstraintError(err) {
			continue
		} else if err != nil {
			return fmt.Errorf("failed to import security %d: %w", sec.ID, err)
		}

		created++
	}

	db.log.Infof("Imported %d out of %d stocks", created, len(secs))

	return nil
}

func (db *DB) AddSubscription(ctx context.Context, securityID int64, userID string) error {
	exist, err := db.Client.Watch.Query().
		Where(
			watch.UserID(userID),
			watch.HasSecurityWith(security.ID(securityID)),
		).
		Exist(ctx)
	if exist {
		return nil
	} else if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("check if watch already exist: %w", err)
	}

	err = db.Client.Watch.
		Create().
		SetUserID(userID).
		SetSecurityID(securityID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("create watch for user %q on security %d: %w", userID, securityID, err)
	}

	return nil
}

func (db *DB) RemoveSubscription(ctx context.Context, securityID int64, userID string) error {
	w, err := db.Client.Watch.Query().Where(
		watch.UserID(userID),
		watch.HasSecurityWith(security.ID(securityID)),
	).First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("fetch watch to remove: %w", err)
	} else if ent.IsNotFound(err) {
		return stockybot.ErrNotFound
	}

	if err = db.Client.Watch.DeleteOne(w).Exec(ctx); err != nil {
		return fmt.Errorf("delete watch: %w", err)
	}

	return nil
}

func (db *DB) GetSubscribedSecurities(
	ctx context.Context,
	userID string,
) ([]stockybot.Security, error) {
	watching, err := db.Client.Watch.
		Query().
		Where(watch.UserID(userID)).
		WithSecurity().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch watching: %w", err)
	}

	watchedSecurities := make([]stockybot.Security, 0, len(watching))

	for _, w := range watching {
		sec := w.Edges.Security

		watchedSecurities = append(watchedSecurities, stockybot.Security{
			ID:       sec.ID,
			Name:     sec.Name,
			List:     sec.List,
			LinkName: sec.LinkName,
			Type:     sec.Type.String(),
			Country:  sec.Country,
		})
	}

	return watchedSecurities, nil
}

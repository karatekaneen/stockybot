package db

import (
	"context"
	"fmt"

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
	client *ent.Client
	log    *zap.SugaredLogger
}

func (db *DB) ImportSecurities(ctx context.Context, provider SecurityProvider) error {
	secs, err := provider.Securities(ctx)
	if err != nil {
		return fmt.Errorf("fetch securities from provider: %w", err)
	}

	db.log.Infof("Found %d stocks to import", len(secs))

	created := 0

	for _, sec := range secs {
		err := db.client.Security.Create().
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
	exist, err := db.client.Watch.Query().
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

	err = db.client.Watch.
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
	w, err := db.client.Watch.Query().Where(
		watch.UserID(userID),
		watch.HasSecurityWith(security.ID(securityID)),
	).First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("fetch watch to remove: %w", err)
	} else if ent.IsNotFound(err) {
		return stockybot.ErrNotFound
	}

	if err = db.client.Watch.DeleteOne(w).Exec(ctx); err != nil {
		return fmt.Errorf("delete watch: %w", err)
	}

	return nil
}

func (db *DB) GetSubscribedSecurities(
	ctx context.Context,
	userID string,
) ([]stockybot.Security, error) {
	watching, err := db.client.Watch.
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

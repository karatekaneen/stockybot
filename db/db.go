package db

import (
	"context"
	"fmt"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/ent"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

type DB struct {
	client *ent.Client
}

func (db *DB) AddSubscription(ctx context.Context, securityID int64, userID string) error {
	exist, err := db.client.Watch.Query().
		Where(
			watch.UserID(userID),
			watch.HasSecurityWith(security.ID(int(securityID))),
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
		SetSecurityID(int(securityID)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("create watch for user %q on security %d: %w", userID, securityID, err)
	}

	return nil
}

func (db *DB) RemoveSubscription(ctx context.Context, securityID int64, userID string) error {
	w, err := db.client.Watch.Query().Where(
		watch.UserID(userID),
		watch.HasSecurityWith(security.ID(int(securityID))),
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
			ID:       int64(sec.ID),
			Name:     sec.Name,
			List:     sec.List,
			LinkName: sec.LinkName,
			Type:     sec.Type.String(),
			Country:  sec.Country,
		})
	}

	return watchedSecurities, nil
}

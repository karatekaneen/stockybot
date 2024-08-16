package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/ent"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

type subscriptionRepository interface{}

type DB struct {
	client *ent.Client
}

func (db *DB) AddSubscription(ctx context.Context, securityId int64, userId string) error {
	exist, err := db.client.Watch.Query().
		Where(
			watch.UserID(userId),
			watch.HasWatchingWith(security.ID(int(securityId))),
		).
		Exist(ctx)
	if exist {
		return nil
	} else if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("check if watch already exist: %w", err)
	}

	err = db.client.Watch.
		Create().
		SetUserID(userId).
		SetWatchingID(int(securityId)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("create watch for user %q on security %d: %w", userId, securityId, err)
	}

	return nil
}

func (db *DB) RemoveSubscription(ctx context.Context, securityId int64, userId string) error {
	return errors.New("not implemented")
}

func (db *DB) GetSubscribedSecurities(
	ctx context.Context,
	userId string,
) ([]stockybot.Security, error) {
	return nil, errors.New("not implemented")
}

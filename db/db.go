package db

import (
	"context"
	"errors"

	"github.com/karatekaneen/stockybot"
	"github.com/karatekaneen/stockybot/ent"
)

type subscriptionRepository interface{}

type DB struct {
	client *ent.Client
}

func (db *DB) AddSubscription(ctx context.Context, securityId int64, userId string) error {
	return errors.New("not implemented")
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

package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/sync/errgroup"
)

// TODO: It would be nice with some tests here. Add when it makes sense to spin up nats on each test

type Store struct {
	stockSubscriberStorage  jetstream.KeyValue // Info about who's subscribing to a particular stock
	userSubscriptionStorage jetstream.KeyValue // Info about users' subscriptions
}

type stockSubscribers struct {
	UserIds []string `json:"users"`
}

type userSubscription struct {
	SecurityIds []int `json:"securities"`
}

func NewStore(ctx context.Context, js jetstream.JetStream) (*Store, error) {
	us, err := js.CreateKeyValue(
		ctx,
		jetstream.KeyValueConfig{Bucket: "user-subscriptions", History: 3},
	)
	if err != nil {
		return nil, fmt.Errorf("user-subscription store init: %w", err)
	}

	ss, err := js.CreateKeyValue(
		ctx,
		jetstream.KeyValueConfig{Bucket: "stock-subscribers", History: 3},
	)
	if err != nil {
		return nil, fmt.Errorf("user-subscription store init: %w", err)
	}

	s := &Store{
		stockSubscriberStorage:  ss,
		userSubscriptionStorage: us,
	}

	return s, nil
}

func (s *Store) AddSubscription(ctx context.Context, securityId int, userId string) error {
	g, ctx := errgroup.WithContext(ctx)

	// Add the subscriber on the stock
	g.Go(func() error { return s.addStockSubscriber(ctx, securityId, userId) })

	// Add the stock to the user's subscribed securities
	g.Go(func() error { return s.addUserSubscription(ctx, securityId, userId) })

	// NOTE: This can become hairy if one fails and one succeeds
	// Is probably not a bad idea to undo on failure in the future
	return g.Wait()
}

func (s *Store) RemoveSubscription(ctx context.Context, securityId int, userId string) error {
	g, ctx := errgroup.WithContext(ctx)

	// Delete user from stock's subscibers
	g.Go(func() error { return s.removeStockSubscriber(ctx, securityId, userId) })

	// // Delete stock from user's subscriptions
	g.Go(func() error { return s.removeUserSubscription(ctx, securityId, userId) })

	// NOTE: This can become hairy if one fails and one succeeds
	// Is probably not a bad idea to undo on failure in the future
	return g.Wait()
}

func (s *Store) GetSubscribedUsers(ctx context.Context, securityId int) ([]string, error) {
	entry, err := s.stockSubscriberStorage.Get(ctx, fmt.Sprint(securityId))
	switch {
	case errors.Is(err, jetstream.ErrKeyNotFound):
		return nil, nil // FIXME: this is not optimal, fix when it's a problem
	case err != nil:
		return nil, fmt.Errorf("get subscribers: %w", err)
	}

	parsed, err := unmarshal[stockSubscribers](entry.Value())
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return parsed.UserIds, nil
}

func (s *Store) GetSubscribedSecurities(ctx context.Context, userId string) ([]int, error) {
	entry, err := s.userSubscriptionStorage.Get(ctx, userId)
	switch {
	case errors.Is(err, jetstream.ErrKeyNotFound):
		return nil, nil // FIXME: this is not optimal, fix when it's a problem
	case err != nil:
		return nil, fmt.Errorf("get user: %w", err)
	}

	parsed, err := unmarshal[userSubscription](entry.Value())
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return parsed.SecurityIds, nil
}

func (s *Store) addUserSubscription(ctx context.Context, securityId int, userId string) error {
	userSubs, err := s.GetSubscribedSecurities(ctx, userId)
	if err != nil {
		return fmt.Errorf("get existing stock subscribers: %w", err)
	}

	userSubs = append(userSubs, securityId)

	b, err := json.Marshal(userSubscription{SecurityIds: userSubs})
	if err != nil {
		return fmt.Errorf("marshal stock subs: %w", err)
	}

	if _, err := s.userSubscriptionStorage.Put(ctx, userId, b); err != nil {
		return fmt.Errorf("put user subscription: %w", err)
	}

	return nil
}

func (s *Store) removeUserSubscription(ctx context.Context, securityId int, userId string) error {
	userSubs, err := s.GetSubscribedSecurities(ctx, userId)
	if err != nil {
		return fmt.Errorf("get existing stock subscribers: %w", err)
	}

	userSubs = slices.DeleteFunc(userSubs, func(id int) bool { return id == securityId })

	b, err := json.Marshal(userSubscription{SecurityIds: userSubs})
	if err != nil {
		return fmt.Errorf("marshal stock subs: %w", err)
	}

	if _, err := s.userSubscriptionStorage.Put(ctx, userId, b); err != nil {
		return fmt.Errorf("put user subscription: %w", err)
	}

	return nil
}

func (s *Store) addStockSubscriber(ctx context.Context, securityId int, userId string) error {
	subs, err := s.GetSubscribedUsers(ctx, securityId)
	if err != nil {
		return fmt.Errorf("get existing stock subscribers: %w", err)
	}

	subs = append(subs, userId)

	b, err := json.Marshal(stockSubscribers{UserIds: subs})
	if err != nil {
		return fmt.Errorf("marshal stock subs: %w", err)
	}

	if _, err := s.stockSubscriberStorage.Put(ctx, fmt.Sprint(securityId), b); err != nil {
		return fmt.Errorf("put stock subscribers: %w", err)
	}

	return nil
}

func (s *Store) removeStockSubscriber(ctx context.Context, securityId int, userId string) error {
	subs, err := s.GetSubscribedUsers(ctx, securityId)
	if err != nil {
		return fmt.Errorf("get existing stock subscribers: %w", err)
	}

	subs = slices.DeleteFunc(subs, func(el string) bool { return el == userId })

	b, err := json.Marshal(stockSubscribers{UserIds: subs})
	if err != nil {
		return fmt.Errorf("marshal stock subs: %w", err)
	}

	if _, err := s.stockSubscriberStorage.Put(ctx, fmt.Sprint(securityId), b); err != nil {
		return fmt.Errorf("put stock subscribers: %w", err)
	}

	return nil
}

func unmarshal[T any](data []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(data, out); err != nil {
		return nil, err
	}
	return out, nil
}

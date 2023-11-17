package bot

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/karatekaneen/stockybot"
)

type subscriptionRepository interface {
	AddSubscription(ctx context.Context, securityId int, userId string) error
	RemoveSubscription(ctx context.Context, securityId int, userId string) error
	GetSubscribedSecurities(ctx context.Context, userId string) ([]int, error)
}

type securityResp struct {
	security *stockybot.Security
	err      error
}

func (b *DiscordBot) watch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	opts := i.ApplicationCommandData().Options
	if len(opts) < 1 {
		b.log.Error("No subcommand specified, this should not happen")
		return
	}

	action := opts[0].Name

	user := getUser(i)

	if action == "list" {
		securities, err := b.listWatchedSecurities(ctx, user)
		if err != nil {
			b.log.Errorw(err.Error(), "user", user.String(), "opts", opts)
			errContent := fmt.Sprintf(
				"An error occured when fetching watched stocks: %s",
				err.Error(),
			)
			if err := interactionResponse(s, i, errContent); err != nil {
				b.log.Error(err)
			}
			return
		}

		secNames := make([]string, 0, len(securities))
		for _, s := range securities {
			secNames = append(secNames, s.Name)
		}

		// Make sure they are alphabetical
		slices.Sort(secNames)

		content := fmt.Sprint("You are watching: %s.", strings.Join(secNames, ", "))
		if err := interactionResponse(s, i, content); err != nil {
			b.log.Error(err)
		}
	}

	optionMap := mapOptions(opts[0].Options)
	ticker := ""
	if opt, ok := optionMap["ticker"]; ok {
		ticker = opt.StringValue()
	}

	if ticker == "" {
		b.log.Errorw("Ticker is empty", "options", optionMap)
		return
	}

	switch action {
	case "add":
		panic("add")
	case "remove":
		panic("remove")
	}
}

func (b *DiscordBot) listWatchedSecurities(
	ctx context.Context,
	user *discordgo.User,
) ([]stockybot.Security, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	secIds, err := b.watchRepo.GetSubscribedSecurities(ctx, user.String())
	if err != nil {
		return nil, fmt.Errorf("could not get subscribed securities: %w", err)
	}

	// Fan out
	ch := make(chan securityResp)
	for _, secId := range secIds {
		go func(id int) {
			sec, err := b.dataRepository.Security(ctx, id)
			if err != nil {
				err = fmt.Errorf("could not fetch security %d: %w", id, err)
			}

			select {
			case <-ctx.Done():
				return
			case ch <- securityResp{security: sec, err: err}:
				// yay!
			}
		}(secId)
	}

	// Fan in
	watchedSecurities := make([]stockybot.Security, 0, len(secIds))
	for range secIds {
		resp := <-ch
		if resp.err != nil {
			return nil, err
		}

		watchedSecurities = append(watchedSecurities, *resp.security)
	}

	return watchedSecurities, nil
}

func getUser(i *discordgo.InteractionCreate) *discordgo.User {
	var user *discordgo.User
	if i.Member != nil {
		user = i.Member.User
	} else {
		user = i.User
	}

	return user
}

package bot

import (
	"context"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"

	"github.com/karatekaneen/stockybot"
)

func (b *DiscordBot) Watch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()

	opts := i.ApplicationCommandData().Options
	if len(opts) < 1 {
		b.log.Error("No subcommand specified, this should not happen")
		return
	}

	action := opts[0].Name

	if action == "list" {
		b.listWatchedSecurities(ctx, s, i)
		return
	}

	user := getUser(i)

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
		// FIXME: Add ticker as int64 here
		if err := b.watchRepo.AddSubscription(ctx, 0, user.String()); err != nil {
			b.log.Errorw("add subscription", "error", err.Error(), "user", user.String())

			errContent := "An error occurred when adding watch: %s" + err.Error()
			if err := interactionResponse(s, i, errContent); err != nil {
				b.log.Error(err)
			}

			return
		}
	case "remove":
		// FIXME: Add ticker as int64 here
		if err := b.watchRepo.RemoveSubscription(ctx, 0, user.String()); err != nil {
			b.log.Errorw("remove subscription", "error", err.Error(), "user", user.String())

			errContent := "An error occurred when remove watch: %s" + err.Error()
			if err := interactionResponse(s, i, errContent); err != nil {
				b.log.Error(err)
			}

			return
		}
	}

	if err := interactionResponse(s, i, "OK"); err != nil {
		b.log.Error(err)
	}
}

func (b *DiscordBot) listWatchedSecurities(
	ctx context.Context,
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	user := getUser(i)

	securities, err := b.getWatchedSecurities(ctx, user)
	if err != nil {
		b.log.Errorw("get watched securities", "error", err.Error(), "user", user.String())

		errContent := "An error occurred when fetching watched stocks: %s" + err.Error()
		if err := interactionResponse(s, i, errContent); err != nil {
			b.log.Error(err)
		}

		return
	}

	content := "You are not watching any stocks."
	if len(securities) > 0 {
		content = "You are watching:\n"
	}

	secNames := make([]string, 0, len(securities))
	for _, s := range securities {
		secNames = append(secNames, s.Name)
	}

	// Make sure they are alphabetical
	slices.Sort(secNames)

	for _, sec := range secNames {
		content = fmt.Sprintf("%s- %s\n", content, sec)
	}

	if err := interactionResponse(s, i, content); err != nil {
		b.log.Error(err)
	}
}

func (b *DiscordBot) getWatchedSecurities(
	ctx context.Context,
	user *discordgo.User,
) ([]stockybot.Security, error) {
	watchedSecs, err := b.watchRepo.GetSubscribedSecurities(ctx, user.String())
	if err != nil {
		return nil, fmt.Errorf("get subscribed securities: %w", err)
	}

	return watchedSecs, nil
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

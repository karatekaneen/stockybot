package bot

import (
	"context"
	"errors"
	"fmt"
	"slices"

	dscd "github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type watchController struct {
	watchRepo subscriptionRepository
	log       *zap.SugaredLogger
}

func (wc *watchController) List(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx := context.Background()

	securities, err := wc.watchRepo.GetSubscribedSecurities(ctx, getUser(i).String())
	if err != nil {
		return fmt.Errorf("get watched securities: %w", err)
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
		return fmt.Errorf("respond to list watched stocks: %w", err)
	}

	return nil
}

func (wc *watchController) AddCommit(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx := context.Background()
	user := getUser(i)

	// Access options in the order provided by the user.
	optionMap := mapOptions(i.ApplicationCommandData().Options)

	ticker := ""

	if opt := optionMap["ticker"]; opt != nil {
		ticker = opt.StringValue()
	}

	if ticker == "" {
		wc.log.Errorw("Ticker is empty", "options", optionMap)
		return errors.New("no ticker provided")
	}

	if err := wc.watchRepo.AddSubscription(ctx, ticker, user.String()); err != nil {
		return fmt.Errorf("add subscription: %w", err)
	}

	return wrapErr(
		interactionResponse(s, i, fmt.Sprintf("Now watching %q", ticker)),
		"respond to AddCommit: %w",
	)
}

func (wc *watchController) AddAutocomplete(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx := context.Background()
	maxResultLen := 10

	optionMap := mapOptions(i.ApplicationCommandData().Options)

	partialTicker := ""

	if opt := optionMap["ticker"]; opt != nil {
		partialTicker = opt.StringValue()
	}

	stockNames, err := wc.watchRepo.GetAllStockNames(ctx)
	if err != nil {
		return fmt.Errorf("get all stock names: %w", err)
	}

	choices := make([]*dscd.ApplicationCommandOptionChoice, 0, maxResultLen)

	// No search, just grab first N stocks if not filtered by fuzzyfind
	allChoices := stockNames
	if partialTicker != "" {
		allChoices = fuzzyFindNItems(stockNames, partialTicker, maxResultLen)
	}

	for i, item := range allChoices {
		if i == maxResultLen {
			break
		}

		choices = append(choices, applicationCommandOptionChoice(item))
	}

	return wrapErr(autocompleteResponse(s, i, choices), "respond to AddAutocomplete: %w")
}

func (wc *watchController) RemoveCommit(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx := context.Background()
	user := getUser(i)

	// Access options in the order provided by the user.
	optionMap := mapOptions(i.ApplicationCommandData().Options)

	ticker := ""

	if opt := optionMap["ticker"]; opt != nil {
		ticker = opt.StringValue()
	}

	if ticker == "" {
		wc.log.Errorw("Ticker is empty", "options", optionMap)
		return errors.New("no ticker provided")
	}

	if err := wc.watchRepo.RemoveSubscription(ctx, ticker, user.String()); err != nil {
		return fmt.Errorf("remove subscription: %w", err)
	}

	return wrapErr(
		interactionResponse(s, i, fmt.Sprintf("Removed watch of %q", ticker)),
		"respond to RemoveCommit: %w",
	)
}

func (wc *watchController) RemoveAutocomplete(s *dscd.Session, i *dscd.InteractionCreate) error {
	ctx := context.Background()
	maxResultLen := 24

	optionMap := mapOptions(i.ApplicationCommandData().Options)

	partialTicker := ""

	if opt := optionMap["ticker"]; opt != nil {
		partialTicker = opt.StringValue()
	}

	watched, err := wc.watchRepo.GetSubscribedSecurities(ctx, getUser(i).String())
	if err != nil {
		return fmt.Errorf("get watched securities: %w", err)
	}

	stockNames := make([]string, 0, len(watched))

	for _, s := range watched {
		stockNames = append(stockNames, s.Name)
	}

	choices := make([]*dscd.ApplicationCommandOptionChoice, 0, maxResultLen)

	// No search, just grab first N stocks if not filtered by fuzzyfind
	allChoices := stockNames
	if partialTicker != "" {
		allChoices = fuzzyFindNItems(stockNames, partialTicker, maxResultLen)
	}

	for i, item := range allChoices {
		if i == maxResultLen {
			break
		}

		choices = append(choices, applicationCommandOptionChoice(item))
	}

	return wrapErr(autocompleteResponse(s, i, choices), "respond to RemoveAutocomplete: %w")
}

func applicationCommandOptionChoice(item string) *dscd.ApplicationCommandOptionChoice {
	return &dscd.ApplicationCommandOptionChoice{Name: item, Value: item}
}

func getUser(i *dscd.InteractionCreate) *dscd.User {
	var user *dscd.User
	if i.Member != nil {
		user = i.Member.User
	} else {
		user = i.User
	}

	return user
}

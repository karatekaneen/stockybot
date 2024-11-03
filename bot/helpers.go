package bot

import (
	"fmt"
	"sort"

	dscd "github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func interactionErr(s *dscd.Session, i *dscd.InteractionCreate, err error) error {
	return s.InteractionRespond(i.Interaction, &dscd.InteractionResponse{
		Type: dscd.InteractionResponseChannelMessageWithSource,
		Data: &dscd.InteractionResponseData{Content: fmt.Sprintf("An error occured: %v", err)},
	})
}

func wrapErr(err error, wrapper string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(wrapper, err)
}

func mapOptions(
	options []*dscd.ApplicationCommandInteractionDataOption,
) map[string]*dscd.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*dscd.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func followUpResponse(s *dscd.Session, i *dscd.Interaction, content string) error {
	_, err := s.FollowupMessageCreate(i, true, &dscd.WebhookParams{Content: content})
	return err
}

func autocompleteResponse(
	s *dscd.Session,
	i *dscd.InteractionCreate,
	choices []*dscd.ApplicationCommandOptionChoice,
) error {
	return s.InteractionRespond(i.Interaction, &dscd.InteractionResponse{
		Type: dscd.InteractionApplicationCommandAutocompleteResult,
		Data: &dscd.InteractionResponseData{Choices: choices},
	})
}

func interactionResponse(s *dscd.Session, i *dscd.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &dscd.InteractionResponse{
		Type: dscd.InteractionResponseChannelMessageWithSource,
		Data: &dscd.InteractionResponseData{Content: content},
	})
}

func fuzzyFindNItems(items []string, subStr string, num int) []string {
	result := fuzzy.RankFindNormalizedFold(subStr, items)
	sort.Sort(result)

	output := make([]string, 0, num)

	for i, item := range result {
		if i == num {
			break
		}

		output = append(output, item.Target)
	}

	return output
}

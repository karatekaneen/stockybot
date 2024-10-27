package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func interactionErr(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	err error,
) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: "An error occured"},
	})
}

func wrapErr(err error, wrapper string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(wrapper, err)
}

func mapOptions(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	return optionMap
}

func followUpResponse(
	s *discordgo.Session,
	i *discordgo.Interaction,
	content string,
) error {
	_, err := s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{Content: content})
	return err
}

func interactionResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	content string,
) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

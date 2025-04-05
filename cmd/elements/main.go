package main

import (
	"tjan-elements/internal/sheets"
	"tjan-elements/internal/twitch"

	"github.com/rs/zerolog/log"
)

const sheetID = "15Fd7qS7_sctKqBjzgwVYs59V3iI4yNo0Exp_G4iDTcI"

func main() {
	googleSheet, err := sheets.NewGoogleSheet()
	if err != nil {
		log.Error().Msgf("Error creating sheet service: %v", err)
		return
	}
	googleSheet.SetId(sheetID)
	chat := twitch.InitChat()
	go chat.StartChat()

	for event := range chat.Events {
		switch event := event.(type) {
		case twitch.Sub:
			err := googleSheet.AddSub(&event)
			if err != nil {
				log.Error().Msgf("Error adding sub to sheet: %v", err)
			}
			log.Info().Msgf("Sub: %s", event.Name)
		case twitch.SubGift:
			log.Info().Msgf("Sub Gift %dx : %s -> %v", *event.Count, event.Gifter, event.Receivers)
		case twitch.Donation:
			log.Info().Msgf("Donation: %s -  €%.2f", event.Name, event.Amount)
		case twitch.Bits:
			log.Info().Msgf("Bits: %s - %d bits", event.Name, event.Amount)
		}

	}
}

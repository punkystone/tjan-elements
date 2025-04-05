package main

import (
	"tjan-elements/internal/sheets"
	"tjan-elements/internal/twitch"

	"github.com/rs/zerolog/log"
)

const sheetID = "1HPPdHKFYDda4OboNPjYKzQupJDxKNA02utjG8HQuTXo"

func main() {
	googleSheet, err := sheets.NewGoogleSheet()
	if err != nil {
		log.Error().Msgf("Error creating sheet service: %v", err)
		return
	}
	googleSheet.SetID(sheetID)
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
			err := googleSheet.AddSubGift(&event)
			if err != nil {
				log.Error().Msgf("Error adding sub gift to sheet: %v", err)
			}
			log.Info().Msgf("Sub Gift %dx : %s -> %v", *event.Count, event.Gifter, event.Receivers)
		case twitch.Donation:
			err := googleSheet.AddDonation(&event)
			if err != nil {
				log.Error().Msgf("Error adding donation to sheet: %v", err)
			}
			log.Info().Msgf("Donation: %s -  €%.2f", event.Name, event.Amount)
		case twitch.Bits:
			err := googleSheet.AddBits(&event)
			if err != nil {
				log.Error().Msgf("Error adding bits to sheet: %v", err)
			}
			log.Info().Msgf("Bits: %s - %d bits", event.Name, event.Amount)
		}
	}
}

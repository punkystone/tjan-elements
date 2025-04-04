package main

import (
	"fmt"
	"tjan-elements/internal/sheets"

	"github.com/rs/zerolog/log"
)

func main() {
	googleSheet, err := sheets.NewGoogleSheet()
	if err != nil {
		log.Error().Msgf("Error creating sheet service: %v", err)
		return
	}
	id := "19TbkQlChmxRHgE3fWRTpUfCf94iQvsWc6w2ezj7CPTs"
	nextCell, _ := googleSheet.GetValue(id, "E1")
	err = googleSheet.UpdateValue(id, fmt.Sprintf("A%s", nextCell.Values[0][0]), "aaaa")
	if err != nil {
		log.Error().Msgf("Error getting values: %v", err)
		return
	}
}

package sheets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"tjan-elements/internal/twitch"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	subIndex      = "Z3"
	subGiftIndex  = "Z4"
	donationIndex = "Z5"
	bitsIndex     = "Z6"
)

type GoogleSheet struct {
	service *sheets.Service
	id      string
}

func CreateAuthURL() (string, error) {
	config, err := createConfig()
	if err != nil {
		return "", err
	}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

func createConfig() (*oauth2.Config, error) {
	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	config, err := google.ConfigFromJSON(credentials, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func CreateToken(authCode string) error {
	config, err := createConfig()
	if err != nil {
		return err
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return err
	}
	file, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(token)
	if err != nil {
		return err
	}
	return nil
}

func NewGoogleSheet() (*GoogleSheet, error) {
	config, err := createConfig()
	if err != nil {
		return nil, err
	}
	file, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(tok)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background(), tok)
	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &GoogleSheet{
		service: service,
	}, nil
}
func (googleSheet *GoogleSheet) SetID(id string) {
	googleSheet.id = id
}

func (googleSheet *GoogleSheet) AddSub(sub *twitch.Sub) error {
	nextSubCell, err := googleSheet.getValue(subIndex)
	if err != nil {
		return err
	}
	err = googleSheet.updateValue(fmt.Sprintf("B%s", nextSubCell.Values[0][0]), [][]any{{sub.Name}})
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) AddSubGift(subgift *twitch.SubGift) error {
	nextGiftCell, err := googleSheet.getValue(subGiftIndex)
	if err != nil {
		return err
	}
	cellValue := nextGiftCell.Values[0][0]
	err = googleSheet.updateValue(fmt.Sprintf("D%s:E%s", cellValue, cellValue), [][]any{{subgift.Gifter, subgift.Count}})
	if err != nil {
		return err
	}
	receivers := [][]any{}
	for _, receiver := range subgift.Receivers {
		receivers = append(receivers, []any{receiver})
	}
	cellValueString, ok := cellValue.(string)
	if !ok {
		return errors.New("invalid cell value type")
	}
	cellValueInt, err := strconv.Atoi(cellValueString)
	if err != nil {
		return err
	}
	cellRange := fmt.Sprintf("G%s:G%d", cellValue, cellValueInt+len(subgift.Receivers))
	err = googleSheet.updateValue(cellRange, receivers)
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) AddDonation(donation *twitch.Donation) error {
	nextDonationCell, err := googleSheet.getValue(donationIndex)
	if err != nil {
		return err
	}
	cellValue := nextDonationCell.Values[0][0]
	err = googleSheet.updateValue(fmt.Sprintf("I%s:J%s", cellValue, cellValue), [][]any{{donation.Name, donation.Amount}})
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) AddBits(bits *twitch.Bits) error {
	nextBitsCell, err := googleSheet.getValue(bitsIndex)
	if err != nil {
		return err
	}
	cellValue := nextBitsCell.Values[0][0]
	err = googleSheet.updateValue(fmt.Sprintf("L%s:M%s", cellValue, cellValue), [][]any{{bits.Name, bits.Amount}})
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) updateValue(field string, values [][]any) error {
	valueRange := &sheets.ValueRange{
		Range:  field,
		Values: values,
	}
	_, err := googleSheet.service.Spreadsheets.Values.Update(googleSheet.id, valueRange.Range, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) getValue(field string) (*sheets.ValueRange, error) {
	response, err := googleSheet.service.Spreadsheets.Values.Get(googleSheet.id, field).Do()
	if err != nil {
		return nil, err
	}
	return response, nil
}

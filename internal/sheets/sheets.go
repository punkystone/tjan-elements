package sheets

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheet struct {
	service *sheets.Service
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

func (googleSheet *GoogleSheet) UpdateValue(id string, field string, value string) error {
	valueRange := &sheets.ValueRange{
		Range:  field,
		Values: [][]any{{value}},
	}
	_, err := googleSheet.service.Spreadsheets.Values.Update(id, valueRange.Range, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return err
	}
	return nil
}

func (googleSheet *GoogleSheet) GetValue(id string, field string) (*sheets.ValueRange, error) {
	response, err := googleSheet.service.Spreadsheets.Values.Get(id, field).Do()
	if err != nil {
		return nil, err
	}
	return response, nil
}

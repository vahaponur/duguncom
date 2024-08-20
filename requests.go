package duguncom

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

const GetOfferURL = "https://api.dugun.com/leads"

type GetOfferParams struct {
	Limit string `json:"limit"`
	Start string `json:"start"`
	End   string `json:"end"`
	Page  string `json:"page"`
}
type Meta struct {
	Total       int `json:"total"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
}
type GetOfferResponse struct {
	Data []Lead `json:"data"`
	Meta Meta   `json:"meta"`
}

func GetOfferRequest(login *LoginResponse, params GetOfferParams) (GetOfferResponse, error) {
	client := resty.New()
	var response GetOfferResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-access-token", login.AccessToken.ID).
		SetHeader("x-consumer-key", login.ConsumerKey).
		SetResult(&response).
		Get(GetOfferURL + "?" + "createdAtStart=" + params.Start + "&createdAtEnd=" + params.End + "&page=" + params.Page + "&limit=" + params.Limit + "&scopes[]=withCoupleTrackingFlag")

	if err != nil {
		return response, err
	}
	if resp.IsError() {
		return response, fmt.Errorf("error: %s", resp.Status())
	}

	return response, nil
}

// SendSmsToCustomers send sms to leads via Dugun.com wait time is milliseconds between requests.
func SendSmsToCustomers(login *LoginResponse, leadIds []string, message string, waitTime int) error {
	client := resty.New()
	baseUrl := "https://api.dugun.com/leads/"
	//IDyi ortaya alan salağa selam olsun
	tail := "/messages/sms"
	type MessageBody struct {
		Body string `json:"body"`
	}
	messageBody := MessageBody{
		Body: message,
	}
	for i, leadId := range leadIds {
		if i == 0 {
			client.Debug = true
		} else {
			client.Debug = false
		}
		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetHeader("x-access-token", login.AccessToken.ID).
			SetHeader("x-consumer-key", login.ConsumerKey).SetBody(messageBody).Post(baseUrl + leadId + tail)
		if err != nil {
			return err
		}
		if resp.IsError() {
			return fmt.Errorf("error: %s", resp.Status())
		}
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}
	return nil

}

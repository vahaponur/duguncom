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

func GetOfferRequest(login LoginResponse, params GetOfferParams) (GetOfferResponse, error) {
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
func SendSmsToCustomers(login LoginResponse, leadIds []string, message string, waitTime int) error {
	client := resty.New()
	baseUrl := "https://api.dugun.com/leads/"
	tail := "/messages/sms"
	type MessageBody struct {
		Body string `json:"body"`
	}
	messageBody := MessageBody{
		Body: message,
	}

	var errList []error

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
			errList = append(errList, fmt.Errorf("error for leadId %s: %w", leadId, err))
			continue // Continue with next leadId
		}

		if resp.IsError() {
			errList = append(errList, fmt.Errorf("error for leadId %s: status %s", leadId, resp.Status()))
			// Continue with next leadId
		}

		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}

	// If there were any errors, return them combined
	if len(errList) > 0 {
		return fmt.Errorf("encountered %d errors: %v", len(errList), errList)
	}

	return nil
}

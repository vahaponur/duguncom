package duguncom

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
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

// SendEmailToCustomers lead'lere email gönderir. waitTime, istekler arasında milisaniye cinsinden bekleme süresidir.
func SendEmailToCustomers(login LoginResponse, leadIds []string, subject string, body string, waitTime int) error {
	client := resty.New()
	baseUrl := "https://api.dugun.com/leads/"
	tail := "/messages"

	type EmailBody struct {
		Body                 string   `json:"body"`
		CarbonCopies         []string `json:"carbonCopies"`
		MessageAttachmentIds []string `json:"messageAttachmentIds"`
		Subject              string   `json:"subject"`
	}

	payload := EmailBody{
		Body:    body,
		Subject: subject,
	}

	var errList []error

	for i, leadId := range leadIds {
		// İlk istekte debug açık kalsın, diğerlerinde kapansın
		if i == 0 {
			client.Debug = true
		} else {
			client.Debug = false
		}

		resp, err := client.R().SetHeader("Content-Type", "application/json").
			SetHeader("x-access-token", login.AccessToken.ID).
			SetHeader("x-consumer-key", login.ConsumerKey).
			SetBody(payload).
			Post(baseUrl + leadId + tail)

		if err != nil {
			errList = append(errList, fmt.Errorf("error for leadId %s: %w", leadId, err))
			continue
		}

		if resp.IsError() {
			errList = append(errList, fmt.Errorf("error for leadId %s: status %s", leadId, resp.Status()))
		}

		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}

	if len(errList) > 0 {
		return fmt.Errorf("encountered %d errors: %v", len(errList), errList)
	}

	return nil
}

// SendSmsOrEmailToCustomers: lead'lerin telefon bilgisine göre SMS veya Email gönderir.
// Telefon numarası boş ya da "null" ise email gönderilir, aksi halde SMS gönderilir.
func SendSmsOrEmailToCustomers(login LoginResponse, leads []Lead, smsMessage string, emailSubject string, emailBody string, waitTime int) error {
	var smsLeadIds []string
	var emailLeadIds []string

	for _, lead := range leads {
		phone := strings.TrimSpace(lead.LeadDetails.Phone)
		if phone == "" || strings.EqualFold(phone, "null") {
			emailLeadIds = append(emailLeadIds, strconv.Itoa(lead.ID))
		} else {
			smsLeadIds = append(smsLeadIds, strconv.Itoa(lead.ID))
		}
	}

	var combinedErr error

	if len(smsLeadIds) > 0 {
		if err := SendSmsToCustomers(login, smsLeadIds, smsMessage, waitTime); err != nil {
			combinedErr = fmt.Errorf("sms errors: %v", err)
		}
	}

	if len(emailLeadIds) > 0 {
		if err := SendEmailToCustomers(login, emailLeadIds, emailSubject, emailBody, waitTime); err != nil {
			if combinedErr != nil {
				combinedErr = fmt.Errorf("%v; email errors: %v", combinedErr, err)
			} else {
				combinedErr = fmt.Errorf("email errors: %v", err)
			}
		}
	}

	return combinedErr
}

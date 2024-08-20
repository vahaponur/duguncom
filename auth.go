package duguncom

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

const LoginURL = "https://api.dugun.com/access-tokens"

type AccessToken struct {
	UserID                 int      `json:"userId"`
	Username               string   `json:"username"`
	ID                     string   `json:"id"`
	UnsignedAgreementIds   []string `json:"unsignedAgreementIds"`
	PendingPayments        []string `json:"pendingPayments"`
	TwoFA                  bool     `json:"2FA"`
	PasswordChangeRequired bool     `json:"passwordChangeRequired"`
	Name                   string   `json:"name"`
}

type LoginResponse struct {
	AccessToken AccessToken `json:"accessToken"`
	ConsumerKey string      `json:"consumerKey"`
	Customer    struct {
		ID int `json:"id"`
	} `json:"customer"`
}

func Login(username, password string, userType ...string) (*LoginResponse, error) {
	client := resty.New()
	utype := "customer"
	if len(userType) > 0 {
		utype = userType[0]
	}
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"username": username,
			"password": password,
			"userType": utype,
		}).
		SetResult(&LoginResponse{}).
		Post(LoginURL)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %s", resp.Status())
	}
	loginResp := resp.Result().(*LoginResponse)
	return loginResp, nil
}

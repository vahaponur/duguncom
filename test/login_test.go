package test

import (
	"duguncom"
	"fmt"
	"strconv"
	"testing"
)

func TestLogin(t *testing.T) {
	resp, err := duguncom.Login("email", "pw")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
}
func TestSms(t *testing.T) {
	resp, err := duguncom.Login("email", "pw")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
	response, err := duguncom.GetOfferRequest(resp, duguncom.GetOfferParams{
		Limit: "100",
		Start: "2024-08-16",
		End:   "2024-08-19",
		Page:  "1",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	leadIds := make([]string, 0)
	for _, data := range response.Data {
		if data.Status == "new" {
			leadIds = append(leadIds, strconv.Itoa(data.ID))
		}
	}
	err = duguncom.SendSmsToCustomers(resp, leadIds, "Merhabalar.\nSizlere düğün.com - Lainvito Design olarak ulaşıyoruz.\nDavetiye modellerimizi incelemek için\nWebsitemiz: https://lainvito.com\nInstagram: @lainvito", 100)
	if err != nil {
		fmt.Println(err.Error())
	}
}

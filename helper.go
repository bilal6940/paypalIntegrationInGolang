package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	constants "test/constant"

	paypalsdk "github.com/gametimesf/paypal-go-sdk"
)

func makePayment(c *paypalsdk.Client, amount paypalsdk.Amount, w http.ResponseWriter, r *http.Request) {
	redirectURI := "http://localhost:9090/success"
	cancelURI := "http://localhost:9090/cancel"
	description := "Description for this payment"

	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
	if err != nil {
		log.Fatalf("error occured in initiating payment --> %v", err)
		return
	}

	fmt.Println("paymentResult", paymentResult)

	for _, v := range paymentResult.Links {
		if v.Rel == "approval_url" {
			fmt.Println("-------url------", v.Href)
			http.Redirect(w, r, v.Href, 302)
		}
	}

	payment, err := c.GetPayment(paymentResult.ID)
	if err != nil {
		log.Fatalf("error occured in getting payment id %v", err)
		return
	}
	fmt.Println("payment", payment.ID)
}

func callRapidAPI(amount paypalsdk.Amount) RapidApiResponse {
	fmt.Println("calling rapid API ------------->")
	to := constants.USD
	from := amount.Currency
	amt := amount.Total
	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/convert?to=%v&from=%v&amount=%v", to, from, amt)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error in making http request %v", err)
	}
	req.Header.Set("apikey", "OwsYwwSOgRBFDkTCKyTgdKrN9C7LQo2j")
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	resBytes, err := ioutil.ReadAll(res.Body)
	rapidAPIResp := RapidApiResponse{}
	if err := json.Unmarshal(resBytes, &rapidAPIResp); err != nil {
		log.Fatalf("error in unmarshilling %v", err)
	}

	fmt.Println(rapidAPIResp)
	return rapidAPIResp
}

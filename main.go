package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	constants "test/constant"

	paypalsdk "github.com/gametimesf/paypal-go-sdk"
	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	var ctx context.Context

	// /pay handler
	mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		//Decoding request body
		reqBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("error in reading request body %v", err)
			return
		}

		reqBody := RequestBody{}
		if err := json.Unmarshal(reqBytes, &reqBody); err != nil {
			log.Fatalf("error in unmarshalling %v", err)
			return
		}

		if err = godotenv.Load(".env"); err != nil {
			log.Fatalf("error in opening env file %v", err)
		}
		clientId := os.Getenv("CLIENT_ID")
		secret := os.Getenv("SECRET")

		c, err := paypalsdk.NewClient(clientId, secret, paypalsdk.APIBaseSandBox)
		if err != nil {
			log.Fatalf("error in creating paypalsdk client %v", err)
			return
		}

		ctx = context.Background()
		ctx = context.WithValue(ctx, "paypalsdkClient", c)

		accessToken, err := c.GetAccessToken()
		if err != nil {
			log.Fatalf("error in accessing token %v", err)
		}
		fmt.Println("token--->", accessToken)

		err = c.SetAccessToken(accessToken.Token)
		if err != nil {
			log.Fatalf("error in setting access token %v", err)
			return
		}

		amount := paypalsdk.Amount{
			Total:    reqBody.Amount,
			Currency: reqBody.Currency,
		}
		var rapidapi RapidApiResponse
		if amount.Currency != constants.USD {
			//call fx API to convert any other currency to USD (base entity for paypal)
			rapidapi = callRapidAPI(amount)
			amount.Total = fmt.Sprintf("%.2f", rapidapi.Info.Rate*rapidapi.Query.Amount)
		}
		amount.Currency = constants.USD
		makePayment(c, amount, w, r)
	})

	//success
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		payerId := r.URL.Query().Get("PayerID")
		paymentId := r.URL.Query().Get("paymentId")

		fmt.Println("---payerId--", payerId)
		fmt.Println("---paymentId--", paymentId)

		value := ctx.Value("paypalsdkClient")
		v := value.(*paypalsdk.Client)
		v.ExecuteApprovedPayment(paymentId, payerId)
		fmt.Println("----payment approved----")
		w.Write([]byte("success"))
	})

	//
	mux.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("cancelled"))
	})

	log.Fatal(http.ListenAndServe(":9090", mux))
}

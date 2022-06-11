package main


type RequestBody struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}


type RapidApiResponse struct {
	Date       string `json:"date"`
	Historical string `json:"historical"`
	Info       Info   `json:"info"`
	Query      Query  `json:"query"`
}

type Query struct {
	Amount float32 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
}
type Info struct {
	Rate      float32 `json:"rate"`
	TimeStamp int64   `json:"time_stamp"`
}

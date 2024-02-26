package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jakobsym/solservice/model"
)

type Token struct{}

// curl -X POST -d '{"coin_address":"7GCihgDB8fe6KNjn2MYtkzZcRjQy3t9GHdC8uHYmW2hr"}' localhost:3000/token
func (t *Token) GetByCA(w http.ResponseWriter, r *http.Request) {
	// call jupv4 get all details about a token
	// display that information
	fmt.Println("getbyca route")
}

func (t *Token) List(w http.ResponseWriter, r *http.Request) {
	// list all tokens in DB
	// here we call jupv4 api to obtain ALL information regarding tickers
	// allowing to display current price
	fmt.Println("List(token) route")
}

func (t *Token) DeleteByCA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deletebyca route")
}

func (t *Token) UpdateByCA(w http.ResponseWriter, r *http.Request) {
	// not sure if needed
	// change the ticker symbol if changed??
	fmt.Println("updatebyca route")
}

func (t *Token) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CoinAddress string `json:"coin_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// query for CA ticker
	token, err := t.FetchTokenData(body.CoinAddress)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if token == nil {
		fmt.Println("error creating token: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: insert `token` into DB

	res, err := json.Marshal(token)
	if err != nil {
		fmt.Println("error marshaling created token: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (t *Token) FetchTokenData(coinAddress string) (*model.Token, error) {
	var response model.Response
	client := &http.Client{}

	url := "https://price.jup.ag/v4/price"
	url += "?ids=" + coinAddress

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating newrequest(): %w", err)
	}

	// send request
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// read response from response body
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding res.body: %w", err)
	}

	// here we create the token model, sending back to other function
	token := &model.Token{
		CoinAddress:   response.Data[coinAddress].CoinAddress,
		MintSymbol:    response.Data[coinAddress].MintSymbol,
		VsTokenSymbol: response.Data[coinAddress].VsTokenSymbol,
	}
	return token, nil
}

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-chi/chi/v5"
	"github.com/jakobsym/solservice/model"
	"github.com/jakobsym/solservice/repository/token"
	"github.com/joho/godotenv"
)

type TokenHandler struct {
	Repo *token.MySqlRepo
}

/*
	if token, err := t.Token.FindByCA(caParam); err != nil {
		// either fetch from DB if token exists in DB already or get info via jup
	}

*/

// curl -X GET localhost:3000/token/HUdqc5MR5h3FssESabPnQ1GTgTcPvnNudAuLj5J6a9sU
func (t *TokenHandler) GetByCA(w http.ResponseWriter, r *http.Request) {
	caParam := chi.URLParam(r, "ca")

	// Get cur price
	price, err := t.FetchTokenPrice(caParam)
	if err != nil {
		fmt.Println("error with FetchTokenPrice()")
	}
	fmt.Println("price: ", price)

	// get supply
	supply, err := GetTokenSupply(caParam)
	if err != nil {
		fmt.Println("error with GetTokenSupply()")
	}
	formattedSupply := FormatFloat(supply)
	fmt.Println("supply: ", formattedSupply)

	// get MC
	mc := CalcMarketCap(supply, price)
	formattedMc := FormatFloat(mc)
	fmt.Println("mc: ", formattedMc)

	// Get all holders
	holders, err := GetTokenHolders(caParam) // 3-15s response time (awful)
	if err != nil {
		fmt.Println("error with GetTokenHolders()")
	}
	fmt.Println("token holders:", holders)
}

func (t *TokenHandler) List(w http.ResponseWriter, r *http.Request) {
	// list all tokens in DB
	// here we call jupv4 api to obtain ALL information regarding tickers
	// allowing to display current price
	fmt.Println("List(token) route")
}

func (t *TokenHandler) DeleteByCA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deletebyca route")
}

func (t *TokenHandler) UpdateByCA(w http.ResponseWriter, r *http.Request) {
	// not sure if needed
	// change the ticker symbol if changed??
	fmt.Println("updatebyca route")
}

// curl -X POST -d '{"coin_address":"7GCihgDB8fe6KNjn2MYtkzZcRjQy3t9GHdC8uHYmW2hr"}' localhost:3000/token
func (t *TokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CoinAddress string `json:"coin_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Creates token
	token, err := t.FetchTokenSymbol(body.CoinAddress)
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
	err = t.Repo.InsertToken(token)
	if err != nil {
		fmt.Println("error inserting token: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(token)
	if err != nil {
		fmt.Println("error marshaling created token: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (t *TokenHandler) FetchTokenSymbol(coinAddress string) (*model.Token, error) {
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

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding res.body: %w", err)
	}

	token := &model.Token{
		CoinAddress:   response.Data[coinAddress].CoinAddress,
		MintSymbol:    response.Data[coinAddress].MintSymbol,
		VsTokenSymbol: response.Data[coinAddress].VsTokenSymbol,
	}
	return token, nil
}

func (t *TokenHandler) FetchTokenPrice(coinAddress string) (float64, error) {
	// either fetch from DB if token exists in DB already or get info via jup
	var response model.Response
	client := &http.Client{}

	url := "https://price.jup.ag/v4/price"
	url += "?ids=" + coinAddress

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating newrequest(): %w", err)
	}

	// send request
	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// read response from response body
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return 0, fmt.Errorf("error decoding res.body: %w", err)
	}
	return float64(response.Data[coinAddress].Price), nil
}

func GetTokenSupply(coinAddress string) (float64, error) {
	endpoint := rpc.MainNetBeta_RPC
	client := rpc.New(endpoint)

	pubKey := solana.MustPublicKeyFromBase58(coinAddress)
	supply, err := client.GetTokenSupply(context.TODO(), pubKey, rpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("error getting token supply: %w", err)
	}
	res, err := strconv.ParseFloat(supply.Value.Amount, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting to uint: %w", err)
	}
	// prettify output?
	return (res / 1e9), nil
}

func CalcMarketCap(supply, price float64) float64 {
	mc := supply * price
	// TODO: format more to only have .2 decimal
	return mc
}

// TODO: Refactor
// Response time range is way too high
// not sure if its how the pagination is implemented
// or if the issue is within my code
func GetTokenHolders(coinAddres string) (uint64, error) {
	client := &http.Client{}
	page := 1
	allHolders := make(map[string]struct{})

	for {
		params := model.Params{
			Page:  page,
			Limit: 1000,
			Mint:  coinAddres,
		}

		body := model.Body{
			JsonRPC: "2.0",
			Method:  "getTokenAccounts",
			Id:      "helius-test",
			Params:  params,
		}

		resBody, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("error marshaling `body`: %w", err)
		}

		url := "https://mainnet.helius-rpc.com/?api-key="
		err = godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
			return 0, fmt.Errorf("error loading data from .env file: %w", err)
		}

		url += os.Getenv("HELIUS_API_KEY")

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(resBody))
		if err != nil {
			return 0, fmt.Errorf("error making post req: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return 0, fmt.Errorf("error sending GetTokenHolders request: %w", err)
		}
		defer res.Body.Close()
		// capture response
		var responseBody model.HolderResponse

		if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			return 0, fmt.Errorf("error decoding res.body: %w", err)
		}
		if responseBody.Result.TokenHolders == nil || len(responseBody.Result.TokenHolders) == 0 {
			break
		}
		// extract token holders to count
		for _, acc := range responseBody.Result.TokenHolders {
			allHolders[acc.Owner] = struct{}{}
		}
		page++
	}

	return uint64(len(allHolders)), nil
}

func FormatFloat(f float64) string {
	formattedFloat := strconv.FormatFloat(f, 'f', -1, 64)
	// TODO: format more to only have .2 decimal
	return formattedFloat
}

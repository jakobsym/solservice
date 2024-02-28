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
	"github.com/jakobsym/solservice/model"
	"github.com/joho/godotenv"
)

type Token struct{}

// localhost:3000/token/<CA>
func (t *Token) GetByCA(w http.ResponseWriter, r *http.Request) {
	//caParam := chi.URLParam(r, "ca")
	/*
		if token, err := t.Token.FindByCA(caParam); err != nil {
			// either fetch from DB if token exists in DB already or get info via jup
		}
	*/
	// show more than just price?
	// do some calculations etc. (% change over X hours)
	// maybe return other misc. information regarding the coin (birdeye api?)
	// I.E: # of holders, MC (Max supply * cur. price), supply?

	/*
		supply: HeliusRPC || Solana RPC API
		holders: HeliusRPC || Solana RPC API (when called, handle error returning server error)
		MC: calulation
	*/

	//fmt.Println("getbyca route")
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

// curl -X POST -d '{"coin_address":"7GCihgDB8fe6KNjn2MYtkzZcRjQy3t9GHdC8uHYmW2hr"}' localhost:3000/token
func (t *Token) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CoinAddress string `json:"coin_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// query for CA ticker
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

func (t *Token) FetchTokenSymbol(coinAddress string) (*model.Token, error) {
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

// TODO: fix
func (t *Token) FetchTokenData(coinAddress string) error {
	// either fetch from DB if token exists in DB already or get info via jup
	var response model.Response
	client := &http.Client{}

	url := "https://price.jup.ag/v4/price"
	url += "?ids=" + coinAddress

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating newrequest(): %w", err)
	}

	// send request
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error doing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// read response from response body
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error decoding res.body: %w", err)
	}
	return nil
}

// ca to test: nosXBVoaCTtYdLvKY6Csb4AC8JCdQKKAaWYtx2ZMoo7
// supply: 99,999,974.58
func GetTokenSupply(coinAddress string) (uint64, error) {
	endpoint := rpc.MainNetBeta_RPC
	client := rpc.New(endpoint)

	pubKey := solana.MustPublicKeyFromBase58(coinAddress)
	supply, err := client.GetTokenSupply(context.TODO(), pubKey, rpc.CommitmentFinalized)
	if err != nil {
		return 0, fmt.Errorf("error getting token supply: %w", err)
	}
	res, err := strconv.ParseUint(supply.Value.Amount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting to uint: %w", err)
	}
	// prettify output?
	return res, nil
}

func CalcMarketCap(supply, price float64) float64 {
	return supply * price
}

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

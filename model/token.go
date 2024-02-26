package model

// get coinaddress from response body
type Token struct {
	CoinAddress   string  `json:"id"`
	MintSymbol    string  `json:"mintSymbol"`
	VsTokenSymbol string  `json:"VsTokenSymbol,omitempty"`
	Price         float64 `json:"price,omitempty"`
}

type Response struct {
	Data map[string]Token `json:"data"`
}

/*
user enters ca which will then display that coins data(price, %change)

user can save a token to a watch list by entering the CA into response body
(get ca, then convert to correct symbol storing the symbol + ca in watchlist)
*/

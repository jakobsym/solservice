package model

// get coinaddress from response body
type Token struct {
	CoinAddress   string  `json:"coin_address"`
	Price         float64 `json:"price"`
	MintSymbol    string  `json:"mint_synbol"`
	VsTokenSymbol string  `json:"vs_token_symbol"`
}

/*

user enters ca which will then display that coins data(price, %change)

user can save a token to a watch list by entering the CA into response body
(get ca, then convert to correct symbol storing the symbol + ca in watchlist)

*/

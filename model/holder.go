package model

type Body struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Id      string `json:"id"`
	Params  Params `json:"params"`
}

type Params struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Mint  string `json:"mint"`
	//DisplayOptions DisplayOptions `json:",omitempty"`
}

type HolderResponse struct {
	JsonRPC string             `json:"jsonrpc"`
	Result  TokenHolderReponse `json:"result"`
}

type TokenHolderReponse struct {
	TokenHolders []TokenHolder `json:"token_accounts"`
}

type TokenHolder struct {
	Owner string `json:"owner"`
}

package token

import (
	"database/sql"
	"fmt"

	"github.com/jakobsym/solservice/model"
)

// database of all tokens I want to add using the Token struct as the attributes

/*
	type Token struct {
		CoinAddress   string  `json:"id"`
		MintSymbol    string  `json:"mintSymbol"`
		VsTokenSymbol string  `json:"VsTokenSymbol,omitempty"`
		Price         float64 `json:"price,omitempty"`
	}
*/
type MySqlRepo struct {
	DB *sql.DB
}

// insert token
// inserts CA + Ticker, returning the ticker to user letting them know what was inserted.
func (ms *MySqlRepo) InsertToken(token *model.Token) error {
	_, err := ms.DB.Exec("INSERT INTO Token(MintSymbol, CoinAddress) VALUES(?,?)", token.MintSymbol, token.CoinAddress)
	if err != nil {
		return fmt.Errorf("error inserting token into db: %w", err)
	}
	return nil
}

//res, err := db.Exec()

// delete token

// search for token

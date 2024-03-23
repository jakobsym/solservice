package token

import (
	"database/sql"
	"fmt"

	"github.com/jakobsym/solservice/model"
)

type MySqlRepo struct {
	DB *sql.DB
}

// insert token
func (ms *MySqlRepo) InsertToken(token *model.Token) error {
	_, err := ms.DB.Exec("INSERT INTO Token(MintSymbol, CoinAddress) VALUES(?,?)", token.MintSymbol, token.CoinAddress)
	if err != nil {
		return fmt.Errorf("error inserting token into db: %w", err)
	}
	return nil
}

// delete token
func (ms *MySqlRepo) DeleteToken(ca string) error {
	_, err := ms.DB.Exec("DELETE FROM Token WHERE CoinAddress = ?" + ca)
	if err != nil {
		return fmt.Errorf("error inserting token into db: %w", err)
	}
	return nil
}

// search for a token
func (ms *MySqlRepo) QueryToken(ca string) (model.Token, error) {
	var token model.Token
	res, err := ms.DB.Query("SELECT CoinAddress FROM Token WHERE CoinAddress = ?" + ca)
	if err != nil {
		return model.Token{}, nil
	}
	if err := res.Scan(&token.CoinAddress, &token.MintSymbol); err != nil {
		return model.Token{}, nil
	}
	return token, nil
}

// search for all tokens
func (ms *MySqlRepo) QueryAllTokens() []model.Token {
	var tokens []model.Token
	rows, err := ms.DB.Query("SELECT * FROM Token.CoinAddress")
	if err != nil {
		return []model.Token{}
	}
	defer rows.Close()
	for rows.Next() {
		var token model.Token
		if err := rows.Scan(&token.CoinAddress, &token.MintSymbol); err != nil {
			return nil
		}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return tokens
}

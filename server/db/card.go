package db

import (
	"database/sql"
)

// SaveCard saves a card to the database.
func (tx Tx) SaveCard(card, name, email string) {
	var (
		oname  string
		oemail string
		err    error
	)
	switch err = tx.tx.QueryRow(`SELECT name, email FROM card_email WHERE card=?`, card).Scan(&oname, &oemail); err {
	case nil, sql.ErrNoRows:
		break
	default:
		panic(err)
	}
	if name == "" && oname != "" {
		name = oname
	}
	if email == "" && oemail != "" {
		email = oemail
	}
	if name == "" && email == "" {
		return
	}
	panicOnNoRows(tx.tx.Exec(`INSERT OR REPLACE INTO card_email (card, name, email) VALUES (?,?,?)`, card, name, email))
}

// FetchCard returns the name and email associated with the card, or empty
// strings if there are none.
func (tx Tx) FetchCard(card string) (name, email string) {
	var (
		err error
	)
	switch err = tx.tx.QueryRow(`SELECT name, email FROM card_email WHERE card=?`, card).Scan(&name, &email); err {
	case nil, sql.ErrNoRows:
		return
	default:
		panic(err)
	}
}

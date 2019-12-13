package db

import (
	"database/sql"

	"scholacantorum.org/orders/model"
)

// SaveCard saves a card to the database.
func (tx Tx) SaveCard(card *model.Card) {
	var (
		oname  string
		oemail string
		err    error
	)
	switch err = tx.tx.QueryRow(`SELECT name, email FROM card_email WHERE card=?`, card.Card).Scan(&oname, &oemail); err {
	case nil, sql.ErrNoRows:
		break
	default:
		panic(err)
	}
	if card.Name == "" && oname != "" {
		card.Name = oname
	}
	if card.Email == "" && oemail != "" {
		card.Email = oemail
	}
	if card.Name == "" && card.Email == "" {
		return
	}
	panicOnNoRows(tx.tx.Exec(`INSERT OR REPLACE INTO card_email (card, name, email) VALUES (?,?,?)`,
		card.Card, card.Name, card.Email))
	tx.audit(model.AuditRecord{Card: card})
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

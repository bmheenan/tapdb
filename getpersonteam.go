package tapdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyGetPersonteam = "getpersonteam"
const qryGetPersonteam = `
SELECT
	email,
	name,
	abbrev,
	colorf,
	colorb
FROM
	personteams
WHERE
	email = ?`

func (db *mySQLDB) initGetPersonteam() error {
	var err error
	db.stmts[keyGetPersonteam], err = db.conn.Prepare(qryGetPersonteam)
	return err
}

// GetPersonteam takes an email and returns a struct with that person or team's details
// Returns a wrapped ErrNotFound if the email doesn't match anything
func (db *mySQLDB) GetPersonteam(email string) (*tapstruct.Personteam, error) {
	if len(email) == 0 {
		return &tapstruct.Personteam{}, errors.New("Email cannot be blank")
	}
	result := db.stmts[keyGetPersonteam].QueryRow(email)
	pt := tapstruct.Personteam{}
	err := result.Scan(&pt.Email, &pt.Name, &pt.Abbrev, &pt.ColorF, &pt.ColorB)
	if errors.Is(err, sql.ErrNoRows) {
		return &tapstruct.Personteam{}, fmt.Errorf("No personteam with that email: %w", ErrNotFound)
	}
	if err != nil {
		return &tapstruct.Personteam{}, fmt.Errorf("Could not get personteam: %v", err)
	}
	return &pt, nil
}

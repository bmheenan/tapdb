package tapdb

import (
	"errors"

	"github.com/bmheenan/tapstruct"
)

const keyGetPersonteam = "getpersonteam"
const qryGetPersonteam = `
`

func (db *mySQLDB) initGetPersonteam() error {
	var err error
	db.stmts[keyGetPersonteam], err = db.conn.Prepare(qryGetPersonteam)
	return err
}

// GetPersonteam takes an email and returns a struct with its details
func (db *mySQLDB) GetPersonteam(email string) (*tapstruct.Personteam, error) {
	return &tapstruct.Personteam{}, errors.New("Not implemented")
}

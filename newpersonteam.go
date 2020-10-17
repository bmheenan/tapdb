package tapdb

import (
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyNewPersonteam = "newpersonteam"
const qryNewPersonteam = `
INSERT INTO personteams (
	email,
	name,
	abbrev,
	colorf,
	colorb
) VALUES (
	?,
	?,
	?,
	?,
	?
)
`

func (db *mySQLDB) initNewPersonteam() error {
	var err error
	db.stmts[keyNewPersonteam], err = db.conn.Prepare(qryNewPersonteam)
	return err
}

func (db *mySQLDB) NewPersonteam(pt *tapstruct.Personteam) error {
	_, err := db.stmts[keyNewPersonteam].Exec(pt.Email, pt.Name, pt.Abbrev, pt.ColorF, pt.ColorB)
	if err != nil {
		return fmt.Errorf("Could not insert new personteam: %v", err)
	}
	return nil
}

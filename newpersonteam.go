package tapdb

import (
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyNewPersonteam = "newpersonteam"
const qryNewPersonteam = `
INSERT INTO personteams (
	email,
	domain,
	name,
	abbrev,
	colorf,
	colorb,
	haschildren
) VALUES (
	?,
	?,
	?,
	?,
	?,
	?,
	FALSE
);`
const keyNewPersonteamParentLink = "newpersonteamparentlink"
const qryNewPersonteamParentLink = `
INSERT INTO personteams_parent_child (
	parent,
	child
) VALUES (
	?,
	?
);`
const keyNewPersonteamUpdateParent = "newpersonteamupdateparent"
const qryNewPersonteamUpdateParent = `
UPDATE personteams
SET
	haschildren = TRUE
WHERE
	email = ?`

func (db *mySQLDB) initNewPersonteam() error {
	var err error
	db.stmts[keyNewPersonteam], err = db.conn.Prepare(qryNewPersonteam)
	if err != nil {
		return err
	}
	db.stmts[keyNewPersonteamParentLink], err = db.conn.Prepare(qryNewPersonteamParentLink)
	if err != nil {
		return err
	}
	db.stmts[keyNewPersonteamUpdateParent], err = db.conn.Prepare(qryNewPersonteamUpdateParent)
	return err
}

// NewPersonteam inserts a new Personteam into the db, with the provided information. It will be a child of the given
// `parentEmail`, or if `parentEmail` == "", it will be inserted at the root of the domain
func (db *mySQLDB) NewPersonteam(pt *tapstruct.Personteam, parentEmail string) error {
	_, err := db.stmts[keyNewPersonteam].Exec(pt.Email, pt.Domain, pt.Name, pt.Abbrev, pt.ColorF, pt.ColorB)
	if err != nil {
		return fmt.Errorf("Could not insert new personteam: %v", err)
	}
	if parentEmail != "" {
		_, errP := db.stmts[keyNewPersonteamParentLink].Exec(parentEmail, pt.Email)
		if errP != nil {
			return fmt.Errorf("Could not link new personteam to parent: %v", errP)
		}
		_, errP = db.stmts[keyNewPersonteamUpdateParent].Exec(parentEmail)
		if errP != nil {
			return fmt.Errorf("Could not update parent's haschildren field: %v", errP)
		}
	}
	return nil
}

package tapdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyGetPersonteam = "getpersonteam"
const qryGetPersonteam = `
SELECT  email,
	    domain,
	    name,
	    abbrev,
	    colorf,
	    colorb,
	    iterationtiming,
	    haschildren
  FROM  personteams
  WHERE email = ?;`
const keyGetPersonteamChildren = "getpersonteamchildren"
const qryGetPersonteamChildren = `
SELECT  child
  FROM  personteams_parent_child
  WHERE parent = ?;`

func (db *mySQLDB) initGetPersonteam() error {
	var err error
	db.stmts[keyGetPersonteam], err = db.conn.Prepare(qryGetPersonteam)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetPersonteam, err)
	}
	db.stmts[keyGetPersonteamChildren], err = db.conn.Prepare(qryGetPersonteamChildren)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetPersonteamChildren, err)
	}
	return nil
}

// GetPersonteam takes an email and returns a struct with that person or team's details. It will also return child
// personteams up to the depth specified. `depth = 0` returns an empty `children` array, `depth = 1` returns max one
// level of children, etc. It may not be larger than 10 (for perf reasons)
// Returns a wrapped ErrNotFound if the provided email doesn't match anything
func (db *mySQLDB) GetPersonteam(email string, depth int) (*tapstruct.Personteam, error) {
	if depth > 10 {
		return &tapstruct.Personteam{}, errors.New("Depth may not be larger than 10")
	}
	pt := &tapstruct.Personteam{}
	err := db.fillInPersonteam(email, depth, pt)
	if err != nil {
		return &tapstruct.Personteam{}, fmt.Errorf("Could not get personteam %v: %w", email, err)
	}
	return pt, nil
}

// This func can be called recursively to fill in a tree of personteams. `pt` is a pointer to the personteam struct to
// write the results into
func (db *mySQLDB) fillInPersonteam(email string, depth int, pt *tapstruct.Personteam) error {
	if len(email) == 0 {
		return errors.New("Email cannot be blank")
	}
	result := db.stmts[keyGetPersonteam].QueryRow(email)
	err := result.Scan(
		&pt.Email,
		&pt.Domain,
		&pt.Name,
		&pt.Abbrev,
		&pt.ColorF,
		&pt.ColorB,
		&pt.IterTiming,
		&pt.HasChildren)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("No personteam with that email: %w", ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("Unexpected error filling in values: %v", err)
	}
	if depth > 0 && pt.HasChildren {
		rows, errChn := db.stmts[keyGetPersonteamChildren].Query(pt.Email)
		if errChn != nil {
			return fmt.Errorf("Could not get children of %v: %v", pt.Email, errChn)
		}
		defer rows.Close()
		for rows.Next() {
			var chEmail string
			errScn := rows.Scan(&chEmail)
			if errChn != nil {
				return fmt.Errorf("Could not scan email from children lookup: %v", errScn)
			}
			child := tapstruct.Personteam{}
			errFill := db.fillInPersonteam(chEmail, depth-1, &child)
			if errFill != nil {
				return fmt.Errorf("Could not fill in details for child %v: %v", chEmail, errFill)
			}
			pt.Children = append(pt.Children, child)
		}
	}
	return nil
}

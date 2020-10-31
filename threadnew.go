package tapdb

import (
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyNewThread = "newthread"
const qryNewThread = `
INSERT INTO threads (name, domain, owner, iteration, state, percentile, costdirect, costtotal)
VALUES              (   ?,      ?,     ?,         ?,     ?,          ?,          ?,         ?);`

func (db *mySQLDB) initThreadNew() error {
	var err error
	db.stmts[keyNewThread], err = db.conn.Prepare(qryNewThread)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyNewThread, err)
	}
	return nil
}

// NewThread inserts a new thread into the db with the given details. You can also link it to existing threads as
// parents or children by providing ids in `parents` or `children`.
// It returns the id of the newly inserted thread if it was inserted
func (db *mySQLDB) ThreadNew(
	name string,
	iteration string,
	owner *tapstruct.Personteam,
	cost int,
	parents []int64,
	children []int64) (int64, error) {

	result, errInsert := db.stmts[keyNewThread].Exec(
		name,
		owner.Domain,
		owner.Email,
		iteration,
		tapstruct.NotStarted,
		1,
		cost,
		cost,
	)
	if errInsert != nil {
		return 0, fmt.Errorf("Could not insert new thread: %v", errInsert)
	}
	id, errID := result.LastInsertId()
	if errID != nil {
		return 0, fmt.Errorf("Could not get new insert id: %v", errID)
	}
	db.StakeholderNew(id, iteration, owner)
	for _, p := range parents {
		errPar := db.ThreadLink(p, id)
		if errPar != nil {
			return id, fmt.Errorf("Could not link new thread to given parent (%v): %v", p, errPar)
		}
	}
	for _, c := range children {
		errCh := db.ThreadLink(id, c)
		if errCh != nil {
			return id, fmt.Errorf("Could not link new thread to given child (%v): %v", c, errCh)
		}
	}
	db.calibrateOrdPct(owner.Email, iteration)
	return id, nil
}

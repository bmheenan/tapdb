package tapdb

import (
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/bmheenan/tapstruct"
)

const keyNewStakeholder = "newstakeholder"
const qryNewStakeholder = `
INSERT INTO threads_stakeholders (
	thread,
	stakeholder,
	domain
  ) VALUES (
	?,
	?,
	?
  );`

func (db *mySQLDB) initNewStakeholder() error {
	var err error
	db.stmts[keyNewStakeholder], err = db.conn.Prepare(qryNewStakeholder)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyNewStakeholder, err)
	}
	return nil
}

// NewStakeholder makes `pt` a stakeholder of the thread with ID `thID`, if not already.
func (db *mySQLDB) NewStakeholder(thID int64, pt *tapstruct.Personteam) error {
	_, errIns := db.stmts[keyNewStakeholder].Exec(thID, pt.Email, pt.Domain)
	if errIns != nil {
		if sqlErr, ok := errIns.(*mysql.MySQLError); ok {
			if sqlErr.Number == 1062 {
				// If it already exists in the database, simply return success
				return nil
			}
		}
		return fmt.Errorf("Could not add new stakeholder: %v", errIns)
	}
	return nil
}

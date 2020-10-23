package tapdb

import (
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyNewThread = "newthread"
const qryNewThread = `
INSERT INTO threads (
	name,
	domain,
	owner,
	iteration,
	state,
	ord,
	percentile,
	costdirect,
	costtotal,
  ) VALUES (
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?,
	?
  );`
const keyNewThreadParentLink = "newthreadparentlink"
const qryNewThreadParentLink = `
INSERT INTO threads_parent_child (
	parent,
	child,
	domain
  ) VALUES (
	?,
	?,
	?
  );`

func (db *mySQLDB) initnewThread() error {
	var err error
	db.stmts[keyNewThread], err = db.conn.Prepare(qryNewThread)
	if err != nil {
		return err
	}
	db.stmts[keyNewThreadParentLink], err = db.conn.Prepare(qryNewThreadParentLink)
	return err
}

// NewThread inserts a new thread into the db with the given details. You can also link it to an existing parent or
// child by providing `parent` or `child` ids. Use 0 for this values to not link to existing threads.
func (db *mySQLDB) NewThread(thread *tapstruct.Threaddetail, parent int64, child int64) (int64, error) {
	_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return 0, fmt.Errorf("Could not `USE` database: %v", errUse)
	}
	result, errInsert := db.stmts[keyNewThread].Exec(
		thread.Name,
		thread.Domain,
		thread.Owner,
		thread.State,
		thread.Order,
		thread.Percentile,
		thread.CostDirect,
		thread.CostTotal)
	if errInsert != nil {
		return 0, fmt.Errorf("Could not insert new thread: %v", errInsert)
	}
	id, errID := result.LastInsertId()
	if errID != nil {
		return 0, fmt.Errorf("Could not get new insert id: %v", errID)
	}
	if parent > 0 {
		_, errParent := db.stmts[keyNewThreadParentLink].Exec(parent, id, thread.Domain)
		if errParent != nil {
			return id, fmt.Errorf("Could not link to given parent thread: %v", errParent)
		}
	}
	if child > 0 {
		_, errChild := db.stmts[keyNewThreadParentLink].Exec(id, child, thread.Domain)
		if errChild != nil {
			return id, fmt.Errorf("Could not link to given child thread: %v", errChild)
		}
	}
	return id, nil
}

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
	costtotal
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

// NewThread inserts a new thread into the db with the given details. You can also link it to existing threads as
// parents or children by providing ids in `parents` or `children`.
// It returns the id of the newly inserted thread if it was inserted
func (db *mySQLDB) NewThread(thread *tapstruct.Threaddetail, parents []int64, children []int64) (int64, error) {
	/*_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return 0, fmt.Errorf("Could not `USE` database: %v", errUse)
	}*/
	result, errInsert := db.stmts[keyNewThread].Exec(
		thread.Name,
		thread.Domain,
		thread.Owner.Email,
		thread.Iteration,
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
	for _, p := range parents {
		_, errParent := db.stmts[keyNewThreadParentLink].Exec(p, id, thread.Domain)
		if errParent != nil {
			return id, fmt.Errorf("Could not link to given parent thread %v: %v", p, errParent)
		}
	}
	for _, c := range children {
		_, errChild := db.stmts[keyNewThreadParentLink].Exec(id, c, thread.Domain)
		if errChild != nil {
			return id, fmt.Errorf("Could not link to given child thread %v: %v", c, errChild)
		}
	}
	return id, nil
}

package tapdb

import (
	"errors"
	"fmt"
	"math"

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
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`
const keyNewThreadParentLink = "newthreadparentlink"
const qryNewThreadParentLink = `
INSERT INTO threads_parent_child (
			  parent,
			  child,
			  domain
			) VALUES (?, ?, ?);`

func (db *mySQLDB) initNewThread() error {
	var err error
	db.stmts[keyNewThread], err = db.conn.Prepare(qryNewThread)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyNewThread, err)
	}
	db.stmts[keyNewThreadParentLink], err = db.conn.Prepare(qryNewThreadParentLink)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyNewThreadParentLink, err)
	}
	return nil
}

// NewThread inserts a new thread into the db with the given details. You can also link it to existing threads as
// parents or children by providing ids in `parents` or `children`. You can specify either, but not both.
// It returns the id of the newly inserted thread if it was inserted
func (db *mySQLDB) NewThread(
	thread *tapstruct.Threaddetail,
	parents []*tapstruct.Threadrow,
	children []*tapstruct.Threadrow) (int64, error) {

	if len(parents) > 0 && len(children) > 0 {
		return 0, errors.New("Cannot insert a new thread with both parents and children. Pick one")
	}
	result, errInsert := db.stmts[keyNewThread].Exec(
		thread.Name,
		thread.Domain,
		thread.Owner.Email,
		thread.Iteration,
		thread.State,
		math.MaxInt32, // temp order
		math.MaxInt32, // temp percentile
		thread.CostDirect,
		thread.CostDirect, // temp costTotal
	)
	if errInsert != nil {
		return 0, fmt.Errorf("Could not insert new thread: %v", errInsert)
	}
	id, errID := result.LastInsertId()
	if errID != nil {
		return 0, fmt.Errorf("Could not get new insert id: %v", errID)
	}
	for _, p := range parents {
		_, errParent := db.stmts[keyNewThreadParentLink].Exec(p.ID, id, thread.Domain)
		if errParent != nil {
			return id, fmt.Errorf("Could not link (ID:%v)to given parent (ID:%v): %v", id, p.ID, errParent)
		}
		if p.Iteration == thread.Iteration {
			db.MoveThread(&tapstruct.Threadrow{
				Owner: tapstruct.Personteam{
					Email: thread.Owner.Email,
				},
				ID:        id,
				Iteration: thread.Iteration,
				Order:     math.MaxInt32,
			}, Before, p)
		}
	}
	if len(parents) == 0 {
		errCal := db.calibrateOrdPct(thread.Owner.Email, thread.Iteration)
		if errCal != nil {
			return id, fmt.Errorf("Could not calibrate iteration after insert: %v", errCal)
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

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
		return err
	}
	db.stmts[keyNewThreadParentLink], err = db.conn.Prepare(qryNewThreadParentLink)
	//if err != nil {
	//	return err
	//}
	//db.stmts[keyGetHighestOrder], err = db.conn.Prepare(qryGetHighestOrder)
	return err
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
		// TODO: Allow NewThread to take both parents and children
		// Complex because you may get parents that are before the children, which means there's no valid order for the
		// new thread. We'll need to move the children up in order or the parents down
	}
	/*var newOrd int
	if len(parents) > 0 {
		// If the new thread has at least one parent, insert it just before the first parent
		ordPar, errOrd := db.conn.Query(fmt.Sprintf(qryGetLowestParentOrder, db.concatInt64AsList(parents)))
		if errOrd != nil {
			return 0, fmt.Errorf("Could not get lowest order of parent threads: %v", errOrd)
		}
		defer ordPar.Close()
		ordMin := 0
		for ordPar.Next() {
			errScn := ordPar.Scan(&ordMin)
			if errScn != nil {
				return 0, fmt.Errorf("Could not scan min parent order: %v", errScn)
			}
		}
		// TODO continue this
	} else if len(children) > 0 {
		// Otherwise, if it has children, insert it right after the last child
	} else {
		// If the new thread has no parents or children, insert it as the end of the iteration
		ord, errOrd := db.stmts[keyGetHighestOrder].Query(thread.Owner.Email, thread.Iteration)
		if errOrd != nil {
			return 0, fmt.Errorf("Could not get the highest order of existing threads: %v", errOrd)
		}
		defer ord.Close()
		highest := 0
		for ord.Next() {
			errScn := ord.Scan(&highest)
			if errScn != nil {
				highest = 0
			}
		}
		newOrd = highest + ((math.MaxInt32 - highest) / 2)
	}*/
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
		_, errParent := db.stmts[keyNewThreadParentLink].Exec(p, id, thread.Domain)
		if errParent != nil {
			return id, fmt.Errorf("Could not link to given parent thread %v: %v", p, errParent)
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
	for _, c := range children {
		_, errChild := db.stmts[keyNewThreadParentLink].Exec(id, c, thread.Domain)
		if errChild != nil {
			return id, fmt.Errorf("Could not link to given child thread %v: %v", c, errChild)
		}
	}
	return id, nil
}

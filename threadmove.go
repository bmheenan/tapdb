package tapdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/bmheenan/tapstruct"
)

const keyGetPrevThreadOrd = "getprevthreadord"
const qryGetPrevThreadOrd = `
SELECT  MAX(ord) AS ord
  FROM  threads
  WHERE owner = ?
    AND iteration = ?
    AND ord < ?;`
const keyGetNextThreadOrd = "getnextthreadord"
const qryGetNextThreadOrd = `
SELECT  MIN(ord) AS ord
  FROM  threads
  WHERE owner = ?
    AND iteration = ?
	AND ord > ?;`
const keyGetPrevThreadPct = "getprevthreadpct"
const qryGetPrevThreadPct = `
SELECT  MAX(ord) AS ord
  FROM  threads
  WHERE owner = ?
	AND iteration = ?
	AND percentile < ?;`
const keyGetNextThreadPct = "getnextthreadpct"
const qryGetNextThreadPct = `
SELECT  MIN(ord) AS ord
  FROM  threads
  WHERE owner = ?
	AND iteration = ?
	AND percentile > ?;`
const keyUpdateOrder = "updateorder"
const qryUpdateOrder = `
UPDATE  threads
  SET   ord = ?
  WHERE id = ?;`

// BeforeAfter specifies if we should put the given thread before or after the reference thread
type BeforeAfter string

const (
	// Before means a lower value for order and percentile
	Before BeforeAfter = "before"
	// After means a higher value for order and percentile
	After BeforeAfter = "after"
)

func (db *mySQLDB) initThreadMove() error {
	var err error
	db.stmts[keyUpdateOrder], err = db.conn.Prepare(qryUpdateOrder)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyUpdateOrder, err)
	}
	db.stmts[keyGetPrevThreadOrd], err = db.conn.Prepare(qryGetPrevThreadOrd)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetPrevThreadOrd, err)
	}
	db.stmts[keyGetNextThreadOrd], err = db.conn.Prepare(qryGetNextThreadOrd)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetNextThreadOrd, err)
	}
	db.stmts[keyGetPrevThreadPct], err = db.conn.Prepare(qryGetPrevThreadPct)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetPrevThreadPct, err)
	}
	db.stmts[keyGetNextThreadPct], err = db.conn.Prepare(qryGetNextThreadPct)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetNextThreadPct, err)
	}
	return nil
}

// MoveThread changes the order and the percentile of the given (first) thread to be either before or after the
// reference (second) thread. `dir` specifies which of those you want: `Before` or `After`. If the owners of the two
// threads are not the same, the move will only be approximate
func (db *mySQLDB) MoveThread(thread *tapstruct.Threadrow, dir BeforeAfter, ref *tapstruct.Threadrow) error {
	// TODO: If `ref` is a different owner than `thread`, the move is approximate. Though it's impossible to make this
	// perfect (e.g. a smaller thread at the very beginning of an iteration will always be before a bigger thread at
	// the beginning of the iteration), it could still be improved by checking if it's moving enough to actually place
	// `thread` before/after `ref`, and if not, trying the next place until it succeeds or gets to the beginning/end
	// of the iteration
	if thread.Iteration != ref.Iteration || thread.Iteration == "" {
		return errors.New("Threads must be in the same (non-empty) iteration")
	}
	if thread.Owner.Email == "" || ref.Owner.Email == "" {
		return errors.New("Threads must have owner emails specified")
	}
	if thread.ID == 0 {
		return errors.New("Thread ID must be specified")
	}
	if thread.Percentile < ref.Percentile && dir == Before {
		// The thread is already before the reference. Nothing needs to be done
		return nil
	}
	if thread.Percentile > ref.Percentile && dir == After {
		// The thread is already after the reference. Nothing needs to be done
		return nil
	}
	var k string
	if dir == Before && thread.Owner.Email == ref.Owner.Email {
		k = keyGetPrevThreadOrd
	} else if dir == After && thread.Owner.Email == ref.Owner.Email {
		k = keyGetNextThreadOrd
	} else if dir == Before && thread.Owner.Email != ref.Owner.Email {
		k = keyGetPrevThreadPct
	} else if dir == After && thread.Owner.Email != ref.Owner.Email {
		k = keyGetNextThreadPct
	}
	var (
		res  *sql.Rows
		errO error
	)
	if thread.Owner.Email == ref.Owner.Email {
		res, errO = db.stmts[k].Query(thread.Owner.Email, thread.Iteration, ref.Order)
	} else {
		res, errO = db.stmts[k].Query(thread.Owner.Email, thread.Iteration, ref.Percentile)
	}
	if errO != nil {
		return fmt.Errorf("Could not query for second thread to insert between: %v", errO)
	}
	defer res.Close()
	var o int
	for res.Next() {
		errScn := res.Scan(&o)
		if errScn != nil && dir == Before {
			// No threads before ref. Move thread to the beginning
			return db.updateThreadOrder(thread.ID, 0, thread.Owner.Email, thread.Iteration)
		}
		if errScn != nil && dir == After {
			// No threads after ref. Move thread to the end
			return db.updateThreadOrder(thread.ID, math.MaxInt32, thread.Owner.Email, thread.Iteration)
		}
	}
	// Move thread in between ref's order and the adjacent thread's order we just fetched from the db
	newOrd := db.min(o, ref.Order) + (db.abs(o-ref.Order) / 2)
	return db.updateThreadOrder(thread.ID, newOrd, thread.Owner.Email, thread.Iteration)
}

func (db *mySQLDB) updateThreadOrder(id int64, ord int, owner string, iter string) error {
	_, errUpd := db.stmts[keyUpdateOrder].Exec(ord, id)
	if errUpd != nil {
		return fmt.Errorf("Could not update order in db: %v", errUpd)
	}
	errCal := db.calibrateOrdPct(owner, iter)
	if errCal != nil {
		return fmt.Errorf("Could not calibrate ord and pct: %v", errCal)
	}
	return nil
}

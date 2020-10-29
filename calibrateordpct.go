package tapdb

import (
	"fmt"
	"math"

	"github.com/bmheenan/tapstruct"
)

const keyGetOrdPct = "getordpct"
const qryGetOrdPct = `
SELECT   id,
		 ord,
		 costdirect,
	     percentile
FROM     threads
WHERE    owner = ?
  AND    iteration = ?
ORDER BY ord;`
const keyGetMaxPctChildren = "getmaxpctchildren"
const qryGetMaxPctChildren = `
SELECT MAX(t.percentile)
FROM   threads t
  JOIN threads_parent_child pc
  ON   t.id = pc.child
WHERE  pc.parent = ?;`
const keyUpdateOrdPct = "updateordpct"
const qryUpdateOrdPct = `
UPDATE threads
SET    ord = ?,
       percentile = ?
WHERE  id = ?;`

func (db *mySQLDB) initCalibrateOrdPct() error {
	var err error
	db.stmts[keyGetOrdPct], err = db.conn.Prepare(qryGetOrdPct)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetOrdPct, err)
	}
	db.stmts[keyGetMaxPctChildren], err = db.conn.Prepare(qryGetMaxPctChildren)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetMaxPctChildren, err)
	}
	db.stmts[keyUpdateOrdPct], err = db.conn.Prepare(qryUpdateOrdPct)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyUpdateOrdPct, err)
	}
	return nil
}

func (db *mySQLDB) calibrateOrdPct(owner string, iter string) error {
	res, errGet := db.stmts[keyGetOrdPct].Query(owner, iter)
	if errGet != nil {
		return fmt.Errorf("Could not query the owner's (%v) iteration (%v): %v", owner, iter, errGet)
	}
	defer res.Close()
	threads := []*tapstruct.Threadrow{}
	ttlCost := 0
	for res.Next() {
		th := tapstruct.Threadrow{}
		errScn := res.Scan(&th.ID, &th.Order, &th.CostCtx, &th.Percentile)
		if errScn != nil {
			return fmt.Errorf("Could not scan threads: %v", errScn)
		}
		ttlCost += th.CostCtx
		threads = append(threads, &th)
	}
	ordStep := math.MaxInt32 / (len(threads) + 1)
	rnCost := 0
	for i, th := range threads {
		newOrd := (i + 1) * ordStep
		rnCost += th.CostCtx
		newPct := float64(rnCost) / float64(ttlCost)
		resCh, errCh := db.stmts[keyGetMaxPctChildren].Query(th.ID)
		if errCh != nil {
			return fmt.Errorf("Could not query children of thread: %v", errCh)
		}
		defer resCh.Close()
		chPct := 0.0
		for resCh.Next() {
			resCh.Scan(&chPct)
			// We ignore the error. A null value means no children, so chPct can be left at 0
		}
		_, errUpd := db.stmts[keyUpdateOrdPct].Exec(newOrd, math.Max(newPct, chPct), th.ID)
		if errUpd != nil {
			return fmt.Errorf("Could not update order and percentile: %v", errUpd)
		}
	}
	return nil
}

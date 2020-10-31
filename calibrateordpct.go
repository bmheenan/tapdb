package tapdb

import (
	"fmt"
	"math"

	"github.com/bmheenan/tapstruct"
)

const keyGetOrdPct = "getordpct"
const qryGetOrdPct = `
SELECT   t.id
  ,      t.costdirect
  ,      t.percentile
  ,      s.ord
FROM     threads t
  JOIN   threads_stakeholders s
  ON     t.id = s.thread
WHERE    t.owner = ?
  AND    s.stakeholder = ?
  AND    t.iteration = ?
ORDER BY ord;`
const keySetOrd = "updateordpct"
const qrySetOrd = `
UPDATE threads_stakeholders
SET    ord = ?
WHERE  thread = ?
  AND  stakeholder = ?;`
const keySetPct = "setpct"
const qrySetPct = `
UPDATE threads
SET    percentile = ?
WHERE  id = ?;`

func (db *mySQLDB) initCalibrateOrdPct() error {
	var err error
	inits := []struct {
		k string
		q string
	}{
		{k: keyGetOrdPct, q: qryGetOrdPct},
		{k: keySetOrd, q: qrySetOrd},
		{k: keySetPct, q: qrySetPct},
	}
	for _, v := range inits {
		db.stmts[v.k], err = db.conn.Prepare(v.q)
		if err != nil {
			return fmt.Errorf("Could not init %v: %v", v.k, err)
		}
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
		errScn := res.Scan(&th.ID, &th.CostCtx, &th.Percentile, &th.Order)
		if errScn != nil {
			return fmt.Errorf("Could not scan threads: %v", errScn)
		}
		ttlCost += th.CostCtx
		threads = append(threads, &th)
	}
	ordStep := math.MaxInt32 / (len(threads) + 1)
	rngCost := 0
	for i, th := range threads {
		newOrd := (i + 1) * ordStep
		rngCost += th.CostCtx
		newPct := float64(rngCost) / float64(ttlCost)
		_, errUOr := db.stmts[keySetOrd].Exec(newOrd, th.ID, owner)
		if errUOr != nil {
			return fmt.Errorf("Could not update order: %v", errUOr)
		}
		_, errUPc := db.stmts[keySetPct].Exec(newPct, th.ID)
		if errUPc != nil {
			return fmt.Errorf("Could not update percentile: %v", errUPc)
		}
	}
	return nil
}

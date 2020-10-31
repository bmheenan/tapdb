package tapdb

import (
	"errors"
	"fmt"

	"github.com/bmheenan/tapstruct"
)

//const keyGetThreadrowsPT = "getthreadrowspt"
const qryGetThreadrowsPT = `
SELECT   t.id
  ,      t.domain
  ,      t.name
  ,      t.owner
  ,      t.iteration
  ,      t.state
  ,      t.percentile
  ,      s.ord
  ,      s.costctx
FROM     threads t
  JOIN   threads_stakeholders s
  ON     t.id = s.thread
WHERE    s.stakeholder = '%v'
  AND    s.toplevel = true
  AND    t.iteration IN (%v)
ORDER BY ord;`

//const keyGetThreadrowChildren = "getthreadrowchildren"
const qryGetThreadrowsChildren = `
SELECT   t.id
  ,      t.domain
  ,      t.name
  ,      t.owner
  ,      t.iteration
  ,      t.state
  ,      t.percentile
  ,      s.ord
  ,      s.costctx
FROM     threads t
  JOIN   threads_stakeholders s
  ON     t.id = s.thread
  JOIN   threads_parent_child pt
  ON     t.id = pt.child
WHERE    pt.parent = %v
ORDER BY ord;`

func (db *mySQLDB) ThreadGetRowsPTPlan(email string, iters []string) ([]tapstruct.Threadrow, error) {
	if email == "" {
		return []tapstruct.Threadrow{}, errors.New("Email cannot be blank")
	}
	if len(iters) == 0 {
		return []tapstruct.Threadrow{}, errors.New("Must include at least one iteration")
	}
	// Get the top level
	il := db.concatStringAsList(iters)
	sqlStmt := fmt.Sprintf(qryGetThreadrowsPT, email, il)
	qRes, errQry := db.conn.Query(sqlStmt)
	if errQry != nil {
		return []tapstruct.Threadrow{}, fmt.Errorf("Could not get list of threads from the db: %v", errQry)
	}
	defer qRes.Close()
	var threads []tapstruct.Threadrow
	for qRes.Next() {
		var (
			th  tapstruct.Threadrow
			oEm string
		)
		errScn := qRes.Scan(
			&th.ID,
			&th.Domain,
			&th.Name,
			&oEm,
			&th.Iteration,
			&th.State,
			&th.Percentile,
			&th.Order,
			&th.CostCtx)
		if errScn != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not scan result from thread query: %v", errScn)
		}
		pt, errPT := db.GetPersonteam(oEm, 0)
		if errPT != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not get owner of thread: %v", errPT)
		}
		th.Owner = *pt
		errPop := db.populateThreadDescendants(&th, oEm, iters)
		if errPop != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf(
				"Could not populate thread descendants of '%v': %v",
				th.Name,
				errPop,
			)
		}
		threads = append(threads, th)
	}
	return threads, nil
}

func (db *mySQLDB) populateThreadDescendants(thread *tapstruct.Threadrow, stakeholder string, iters []string) error {
	sqlStmt := fmt.Sprintf(qryGetThreadrowsChildren, thread.ID)
	chs, errQry := db.conn.Query(sqlStmt)
	if errQry != nil {
		return fmt.Errorf("Could not query children of '%v': %v", thread.Name, errQry)
	}
	chs.Close()
	for chs.Next() {
		var (
			th  tapstruct.Threadrow
			oEm string
		)
		errScn := chs.Scan(
			&th.ID,
			&th.Domain,
			&th.Name,
			&oEm,
			&th.Iteration,
			&th.State,
			&th.Percentile,
			&th.Order,
			&th.CostCtx,
		)
		if errScn != nil {
			return fmt.Errorf("Could not scan child threadrow: %v", errScn)
		}
		pt, errPT := db.GetPersonteam(oEm, 0)
		if errPT != nil {
			return fmt.Errorf("Could not get owner of threadrow: %v", errPT)
		}
		th.Owner = *pt
		db.populateThreadDescendants(&th, stakeholder, iters)
		if db.strIncl(&iters, th.Iteration) && stakeholder == th.Owner.Email {
			thread.Children = append(thread.Children, th)
		} else {
			for _, ch := range th.Children {
				if db.strIncl(&iters, ch.Iteration) && stakeholder == ch.Owner.Email {
					thread.Children = append(thread.Children, ch)
				}
			}
		}
	}
	return nil
}

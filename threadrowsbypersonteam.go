package tapdb

import (
	"errors"
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyGetThreadrowsByPT = "getthreadrowsbypersonteam"
const qryGetThreadrowsByPT = `
SELECT
	# dedupe any rows with a personteam listed as both an owner and stakeholder
	UNIQUE id,
	MAX(domain) AS domain,
	MAX(name) AS name,
	MAX(state) AS state,
	MAX(costdirect) AS costdirect,
	MAX(owner) AS owner,
	MAX(iteration) AS iteration,
	MAX(ord) AS ord,
	MAX(percentile) AS percentile
  FROM
	(
	SELECT
		id,
		domain,
		name,
		state,
		costdirect,
		owner,
		iteration,
		ord,
		percentile
	FROM
		threads
	WHERE
			owner = ?
		AND iteration IN ( ? )
	UNION
	SELECT
		t.id,
		t.domain,
		t.name,
		t.state,
		t.costdirect,
		t.owner,
		t.iteration,
		t.ord,
		t.percentile	
	FROM
		threads AS t
	INNER JOIN
		threads_stakeholders AS s
	ON
		t.id = s.thread
	WHERE
			s.stakeholder = ?
		AND t.iteration IN ( ? )
	) AS unioned_results;`

const keyGetThreadrowsPC = "getthreadrowsparentchild"
const qryGetThreadrowsPC = `
SELECT
	parent,
	child
  FROM
	threads_parent_child
  WHERE
	parent IN ( ? );`

func (db *mySQLDB) initGetThreadrowsByPersonteamPlan() error {
	var err error
	db.stmts[keyGetThreadrowsByPT], err = db.conn.Prepare(qryGetThreadrowsByPT)
	if err != nil {
		return err
	}
	db.stmts[keyGetThreadrowsPC], err = db.conn.Prepare(qryGetThreadrowsPC)
	return err
}

func (db *mySQLDB) GetThreadrowsByPersonteamPlan(email string, iters []string) ([]tapstruct.Threadrow, error) {
	if email == "" {
		return []tapstruct.Threadrow{}, errors.New("Email cannot be blank")
	}
	if len(iters) == 0 {
		return []tapstruct.Threadrow{}, errors.New("Must include at least one iteration")
	}
	_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return []tapstruct.Threadrow{}, fmt.Errorf("Could not `USE` database: %v", errUse)
	}
	il := db.concatIters(iters)
	qRes, errQry := db.stmts[keyGetThreadrowsByPT].Query(email, il, email, il)
	if errQry != nil {
		return []tapstruct.Threadrow{}, fmt.Errorf("Could not get list of threads from the db: %v", errQry)
	}
	defer qRes.Close()
	flatThs := []tapstruct.Threadrow{}
	for qRes.Next() {
		th := tapstruct.Threadrow{}
		errScn := qRes.Scan(
			&th.ID,
			&th.Domain,
			&th.Name,
			&th.State,
			&th.CostCtx,
			&th.Owner,
			&th.Iteration,
			&th.Order,
			&th.Percentile)
		if errScn != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not scan result: %v", errScn)
		}
		flatThs = append(flatThs, th)
	}
	pcRes, errPCQry := db.stmts[keyGetThreadrowsPC].Query(qryGetThreadrowsPC)
	if errQry != nil {
		return []tapstruct.Threadrow{}, fmt.Errorf("Could not query thread parent/child relationships: %v", errPCQry)
	}
	defer pcRes.Close()
	rels := []threadRel{}
	for pcRes.Next() {
		rel := threadRel{}
		errScn := pcRes.Scan(&rel.parent, &rel.child)
		if errScn != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not scan parent/child relationship: %v", errScn)
		}
		rels = append(rels, rel)
	}
	threads := []tapstruct.Threadrow{}
	for _, v := range flatThs {
		if !db.isInChildren(rels, v.ID) {
			th, errTFA := db.threadrowFromArray(flatThs, v.ID)
			if errTFA != nil {
				return []tapstruct.Threadrow{}, fmt.Errorf("Could not find the expected thread: %v", errTFA)
			}
			threads = append(threads, th)
		}
	}
	return threads, nil
}

type threadRel struct {
	parent int
	child  int
}

func (db *mySQLDB) concatIters(iters []string) string {
	res := ""
	l := len(iters) - 1
	for i, v := range iters {
		res = res + fmt.Sprintf("'%s'", v)
		if i < l {
			res = res + ", "
		}
	}
	return res
}

func (db *mySQLDB) isInChildren(rels []threadRel, id int) bool {
	for _, v := range rels {
		if v.child == id {
			return true
		}
	}
	return false
}

func (db *mySQLDB) threadrowFromArray(threads []tapstruct.Threadrow, id int) (tapstruct.Threadrow, error) {
	for _, v := range threads {
		if v.ID == id {
			return v, nil
		}
	}
	return tapstruct.Threadrow{}, errors.New("No thread by that ID found")
}

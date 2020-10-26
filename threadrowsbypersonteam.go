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
	id,
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
			owner = '%v'
		AND iteration IN ( %v )
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
			s.stakeholder = '%v'
		AND t.iteration IN ( %v )
	) AS unioned_results
  GROUP BY
	id;`

const keyGetThreadrowAncestors = "getthreadrowancestors"
const qryGetThreadrowAncestors = `
WITH RECURSIVE ancestors (child, parent) AS
	(
	SELECT
		child,
		parent
	  FROM
		threads_parent_child
	  WHERE
		child = ?
	UNION ALL
	SELECT
		t.child,
		t.parent
	  FROM
		threads_parent_child t
	  INNER JOIN
		ancestors
	  ON
		t.child = ancestors.parent
	)
SELECT
	parent
  FROM
	ancestors
  ORDER BY
	parent;`

/*`
SELECT
	parent,
	child
  FROM
	threads_parent_child
  WHERE
	child = ?;`*/

func (db *mySQLDB) initGetThreadrowsByPersonteamPlan() error {
	var err error
	/*db.stmts[keyGetThreadrowsByPT], err = db.conn.Prepare(qryGetThreadrowsByPT)
	if err != nil {
		return err
	}*/
	db.stmts[keyGetThreadrowAncestors], err = db.conn.Prepare(qryGetThreadrowAncestors)
	return err
}

func (db *mySQLDB) GetThreadrowsByPersonteamPlan(email string, iters []string) ([]tapstruct.Threadrow, error) {
	if email == "" {
		return []tapstruct.Threadrow{}, errors.New("Email cannot be blank")
	}
	if len(iters) == 0 {
		return []tapstruct.Threadrow{}, errors.New("Must include at least one iteration")
	}
	il := db.concatIters(iters)
	sqlStmt := fmt.Sprintf(qryGetThreadrowsByPT, email, il, email, il)
	qRes, errQry := db.conn.Query(sqlStmt)
	if errQry != nil {
		return []tapstruct.Threadrow{}, fmt.Errorf("Could not get list of threads from the db: %v", errQry)
	}
	defer qRes.Close()
	threads := []*threadWMeta{}
	for qRes.Next() {
		th := threadWMeta{
			thread: tapstruct.Threadrow{},
		}
		var oEmail string
		errScn := qRes.Scan(
			&th.thread.ID,
			&th.thread.Domain,
			&th.thread.Name,
			&th.thread.State,
			&th.thread.CostCtx,
			&oEmail,
			&th.thread.Iteration,
			&th.thread.Order,
			&th.thread.Percentile)
		if errScn != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not scan result from thread query: %v", errScn)
		}
		pt, errPT := db.GetPersonteam(oEmail, 0)
		if errPT != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not get owner of thread: %v", errPT)
		}
		th.thread.Owner = *pt
		ancestors, errAn := db.stmts[keyGetThreadrowAncestors].Query(th.thread.ID)
		if errAn != nil {
			return []tapstruct.Threadrow{}, fmt.Errorf("Could not get parents of thread: %v", errAn)
		}
		defer ancestors.Close()
		for ancestors.Next() {
			var a int
			errScn := ancestors.Scan(&a)
			if errScn != nil {
				return []tapstruct.Threadrow{}, fmt.Errorf("Could not scan ancestor of thread: %v", errScn)
			}
			th.parents = append(th.parents, a)
		}
		threads = append(threads, &th)
	}
	threads = db.nestThreads(threads)
	return db.stripMeta(&threads), nil
}

type threadWMeta struct {
	thread   tapstruct.Threadrow
	parents  []int
	children []*threadWMeta
}

func (db *mySQLDB) nestThreads(threads []*threadWMeta) []*threadWMeta {
	for i := 0; i < len(threads); {
		for j := 0; j < len(threads); j++ {
			if db.isInChildren(threads[j].parents, threads[i].thread.ID) {
				threads[j].children = append(threads[j].children, threads[i])
				db.removeInt(&(threads[i].parents), threads[j].thread.ID)
				db.removeThread(&threads, i)
				threads[i].children = db.nestThreads(threads[i].children)
			} else {
				i++
			}
		}
	}
	return threads
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

func (db *mySQLDB) isInChildren(a []int, id int) bool {
	for _, v := range a {
		if v == id {
			return true
		} else if v > id {
			return false
		}
	}
	return false
}

func (db *mySQLDB) removeInt(a *[]int, item int) {
	index := -1
	for i, v := range *a {
		if v == item {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	copy((*a)[index:], (*a)[index+1:])
	(*a)[len(*a)-1] = 0
	*a = (*a)[:len(*a)-1]
}

func (db *mySQLDB) removeThread(a *[]*threadWMeta, i int) {
	copy((*a)[i:], (*a)[i+1:])
	(*a)[len(*a)-1] = nil
	*a = (*a)[:len(*a)-1]
}

func (db *mySQLDB) stripMeta(in *[]*threadWMeta) []tapstruct.Threadrow {
	ret := []tapstruct.Threadrow{}
	for _, v := range *in {
		th := v.thread
		th.Children = db.stripMeta(&v.children)
		ret = append(ret, th)
	}
	return ret
}

/*
type threadRel struct {
	parent int
	child  int
}
*/

/*
func (db *mySQLDB) threadrowFromArray(threads []tapstruct.Threadrow, id int) (tapstruct.Threadrow, error) {
	for _, v := range threads {
		if v.ID == id {
			return v, nil
		}
	}
	return tapstruct.Threadrow{}, errors.New("No thread by that ID found")
}
*/

// TODO: nest threads with parents in the list under their respective parent
//			While at least one thread on the top level has at least one parent:
//            Find the first thread with a parent
//          	If a parent is in the list, nest it there
//              Else if there are parents but none are on the list, and expanded seach is flase, add all ancestors and set expanded search true
//              Else if parents are not on the list, and expanded search is true, delete all parents

/*
	pcRes, errPCQry := db.stmts[keyGetThreadrowsPC].Query()
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
	// Top level of threads should only have those which aren't children of others
	for _, v := range flatThs {
		if !db.isInChildren(rels, v.ID) {
			th, errTFA := db.threadrowFromArray(flatThs, v.ID)
			if errTFA != nil {
				return []tapstruct.Threadrow{}, fmt.Errorf("Could not find the expected thread: %v", errTFA)
			}
			threads = append(threads, th)
		}
	}*/

package tapdb

import (
	"errors"
	"fmt"

	"github.com/bmheenan/tapstruct"
)

const keyGetThreadrowsByPT = "getthreadrowsbypersonteam"
const qryGetThreadrowsByPT = `
SELECT   id,
	     MAX(domain) AS domain,
	     MAX(name) AS name,
	     MAX(state) AS state,
	     MAX(costdirect) AS costdirect,
	     MAX(owner) AS owner,
	     MAX(iteration) AS iteration,
	     MAX(ord) AS ord,
		 MAX(percentile) AS percentile
FROM     (
         SELECT id,
		        domain,
		        name,
		        state,
		        costdirect,
		        owner,
		        iteration,
		        ord,
		        percentile
	     FROM   threads
	     WHERE  owner = '%v'
	       AND  iteration IN ( %v )
	     UNION
	     SELECT t.id,
		        t.domain,
		        t.name,
		        t.state,
		        t.costdirect,
		        t.owner,
		        t.iteration,
		        t.ord,
		        t.percentile	
	     FROM   threads AS t
	       JOIN threads_stakeholders AS s
	       ON   t.id = s.thread
	     WHERE  s.stakeholder = '%v'
		   AND  t.iteration IN ( %v )
	     ) AS   unioned_results
GROUP BY id
ORDER BY percentile,
         ord;`

const keyGetThreadrowAncestors = "getthreadrowancestors"
const qryGetThreadrowAncestors = `
WITH     RECURSIVE ancestors (child, parent) AS
         (
         SELECT child,
                parent
         FROM   threads_parent_child
         WHERE  child = ?
         UNION ALL
         SELECT t.child,
                t.parent
         FROM   threads_parent_child t
         JOIN   ancestors
           ON   t.child = ancestors.parent
	     )
SELECT   parent
FROM     ancestors
ORDER BY parent;`

func (db *mySQLDB) initGetThreadrowsByPersonteamPlan() error {
	var err error
	db.stmts[keyGetThreadrowAncestors], err = db.conn.Prepare(qryGetThreadrowAncestors)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyGetThreadrowAncestors, err)
	}
	return nil
}

func (db *mySQLDB) GetThreadrowsByPersonteamPlan(email string, iters []string) ([]tapstruct.Threadrow, error) {
	if email == "" {
		return []tapstruct.Threadrow{}, errors.New("Email cannot be blank")
	}
	if len(iters) == 0 {
		return []tapstruct.Threadrow{}, errors.New("Must include at least one iteration")
	}
	il := db.concatStringAsList(iters)
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
			var a int64
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
	parents  []int64
	children []*threadWMeta
}

func (db *mySQLDB) nestThreads(threads []*threadWMeta) []*threadWMeta {
	for i := 0; i < len(threads); {
		remove := false
		for j := 0; j < len(threads); j++ {
			if db.isInSortedInt(threads[i].parents, threads[j].thread.ID) &&
				!db.isInThreadrow(threads[j].children, threads[i].thread.ID) {
				threads[j].children = append(threads[j].children, threads[i])
				db.removeInt(&(threads[i].parents), threads[j].thread.ID)
				remove = true
			}
		}
		if remove {
			db.removeThread(&threads, i)
		} else {
			i++
		}
	}
	for _, v := range threads {
		v.children = db.nestThreads(v.children)
	}
	return threads
}

func (db *mySQLDB) isInSortedInt(a []int64, id int64) bool {
	for _, v := range a {
		if v == id {
			return true
		} else if v > id {
			return false
		}
	}
	return false
}

func (db *mySQLDB) isInThreadrow(a []*threadWMeta, id int64) bool {
	for _, v := range a {
		if v.thread.ID == id {
			return true
		}
	}
	return false
}

func (db *mySQLDB) removeInt(a *[]int64, item int64) {
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

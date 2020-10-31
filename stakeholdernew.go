package tapdb

import (
	"fmt"
	"math"

	"github.com/go-sql-driver/mysql"

	"github.com/bmheenan/tapstruct"
)

const keyNewStakeholder = "newstakeholder"
const qryNewStakeholder = `
INSERT INTO threads_stakeholders (thread, stakeholder, domain, ord, toplevel, costctx)
VALUES                           (     ?,           ?,      ?,   ?,        ?,       ?);`
const keyGetThreadAncestors = "getthreadrowancestors"
const qryGetThreadAncestors = `
WITH     RECURSIVE ancestors (child, parent) AS
         (
         SELECT child
           ,    parent
         FROM   threads_parent_child
         WHERE  child = ?
         UNION ALL
         SELECT t.child
           ,    t.parent
         FROM   threads_parent_child t
         JOIN   ancestors
           ON   t.child = ancestors.parent
	     )
SELECT   a.parent
  ,      s.stakeholder
  ,      t.iteration
FROM     ancestors a
  JOIN   threads_stakeholders s
  ON     s.thread = a.parent
  JOIN   threads t
  ON     t.id = a.parent
ORDER BY s.stakeholder;`
const keyGetThreadDescendants = "getthreadrowdescendants"
const qryGetThreadDescendants = `
WITH     RECURSIVE descendants (child, parent) AS
         (
         SELECT child
           ,    parent
         FROM   threads_parent_child
         WHERE  parent = ?
         UNION ALL
         SELECT t.child
           ,    t.parent
         FROM   threads_parent_child t
         JOIN   descendants
           ON   t.parent = descendants.child
	     )
SELECT   d.child
  ,      t.owner
  ,      s.stakeholder
  , 	 t.costdirect
  , 	 t.iteration
FROM     descendants d
  JOIN   threads_stakeholders s
  ON     s.thread = d.child
  JOIN   threads t
  ON     t.id = d.child
ORDER BY s.stakeholder;`
const keySetThreadToplevel = "setthreadtoplevel"
const qrySetThreadToplevel = `
UPDATE threads_stakeholders
SET    toplevel = ?
WHERE  thread = ?
  AND  stakeholder = ?`
const keyGetPTDescendants = "gettptdescendants"
const qryGetPTDescendants = ``

func (db *mySQLDB) initStakeholderNew() error {
	var err error
	inits := []struct {
		k string
		q string
	}{
		{k: keyNewStakeholder, q: qryNewStakeholder},
		{k: keyGetThreadAncestors, q: qryGetThreadAncestors},
		{k: keyGetThreadDescendants, q: qryGetThreadDescendants},
		{k: keySetThreadToplevel, q: qrySetThreadToplevel},
		{k: keyGetPTDescendants, q: qryGetPTDescendants},
	}
	for _, v := range inits {
		db.stmts[v.k], err = db.conn.Prepare(v.q)
		if err != nil {
			return fmt.Errorf("Could not init %v: %v", v.k, err)
		}
	}
	return nil
}

// NewStakeholder makes `pt` a stakeholder of the thread with ID `id`, if not already. The iteration of the thread
// must be provided.
func (db *mySQLDB) StakeholderNew(id int64, iteration string, pt *tapstruct.Personteam) error {
	// Determine if the thread is top level (not a child of another thread with the same stakeholder)
	top := true
	ancs, errAn := db.stmts[keyGetThreadAncestors].Query(id)
	if errAn != nil {
		return fmt.Errorf("Could not query thread ancestors: %v", errAn)
	}
	defer ancs.Close()
	for ancs.Next() {
		var (
			i  int64
			s  string
			it string
		)
		errScn := ancs.Scan(&i, &s, &it)
		if errScn != nil {
			return fmt.Errorf("Could not scan ancestors: %v", errScn)
		}
		if i != id && s == pt.Email && it == iteration {
			top = false
		}
	}
	// Determine the costCtx, the cost of the thread for this stakeholder (including subteams + team members), and
	// if any threads are no longer the top level as a result of this new stakeholder
	cost := 0
	displaced := []int64{}
	mbrs, errM := db.stmts[keyGetPTDescendants].Query(pt.Email)
	if errM != nil {
		return fmt.Errorf("Could not query sub teams + team members: %v", errM)
	}
	defer mbrs.Close()
	teamMembers := []string{}
	for mbrs.Next() {
		var m = ""
		errScn := mbrs.Scan(&m)
		if errScn != nil {
			return fmt.Errorf("Could not scan sub team members: %v", errScn)
			//break
		}
		teamMembers = append(teamMembers, m)
	}
	des, errDe := db.stmts[keyGetThreadDescendants].Query(id)
	if errDe != nil {
		return fmt.Errorf("Could not query thread descendants: %v", errDe)
	}
	defer des.Close()
	for des.Next() {
		var (
			i  int64
			o  string
			s  string
			c  int
			it string
		)
		errScn := des.Scan(&i, &o, &s, &c, &it)
		if errScn != nil {
			return fmt.Errorf("Could not scan from descendants: %v", errScn)
		}
		if i != id && db.ptIncludes(pt, s) && it == iteration {
			displaced = append(displaced, i)
		}
		if db.strIncl(&teamMembers, o) {
			cost += c
		}
	}
	// Make pt a stakeholder of the thread, putting it at the end of their iteration
	_, errIns := db.stmts[keyNewStakeholder].Exec(
		id,
		pt.Email,
		pt.Domain,
		math.MaxInt32,
		top,
		cost,
	)
	if errIns != nil {
		if sqlErr, ok := errIns.(*mysql.MySQLError); ok {
			if sqlErr.Number == 1062 {
				// If it already exists in the database, simply return success, without actually making any changes
				return nil
			}
		}
		return fmt.Errorf("Could not add new stakeholder: %v", errIns)
	}
	// Set threads that are no longer the top level for this stakeholder + iteration as such. They'll show as children
	for _, v := range displaced {
		_, errSet := db.stmts[keySetThreadToplevel].Exec(false, v, pt.Email)
		if errSet != nil {
			return fmt.Errorf("Could not set displaced thread as not top level: %v", errSet)
		}
	}
	return nil
}

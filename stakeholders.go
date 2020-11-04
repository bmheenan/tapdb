package tapdb

import (
	"fmt"

	taps "github.com/bmheenan/taps"
)

// NewStakeholder makes `stakeholder` a stakeholder of `thread`. It also requires that we specify the `domain` that owns
// the relationship, the `iteration` the thread will show up in for that stake holder, the `order` of this thread for
// this stakedholder, whether this thread is the top level of the iteration for this stakeholder (`topLvl`), and the
// `cost` for this stakeholder to complete the thread.
func (db *mysqlDB) NewStakeholder(
	thread int64,
	stakeholder,
	domain,
	iter string,
	order int,
	topLvl bool,
	cost int,
) error {
	if stakeholder == "" || domain == "" || order < 0 || cost < 0 {
		return fmt.Errorf("Stakeholder and domain must be non-blank; Ord and cost must be >= 0: %w", ErrBadArgs)
	}
	_, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders
	            (thread, stakeholder, domain, iteration, ord, toplevel, costctx)
	VALUES      (    %v,        '%v',   '%v',      '%v',  %v,       %v,      %v)
	;`, thread, stakeholder, domain, iter, order, topLvl, cost))
	if errIn != nil {
		return fmt.Errorf("Could not add stakeholder: %v", errIn)
	}
	return nil
}

// GetAffectedStakeholders takes `parent` and `child` thread ids, and returns a map of all personteams that are a
// stakeholder of at least one ancestor of the parent and at least one descendant of the child. They'll be affected by
// (un)linking the two threads.
func (db *mysqlDB) GetStakeholderAncestors(thread int64) (map[string]*taps.Personteam, error) {
	return db.getStakeholderAncDes(`
	WITH   RECURSIVE ancestors (child, parent) AS
	       (
	       SELECT child
			 ,    parent
	       FROM   threads_parent_child
		   WHERE  child = %v
	       UNION ALL
	       SELECT t.child
			 ,    t.parent
	       FROM   threads_parent_child t
	         JOIN ancestors a
			 ON   t.child = a.parent
		   )
	SELECT DISTINCT stakeholder
	FROM   ancestors a
	  JOIN threads_stakeholders t
	  ON   a.parent = t.thread
	;`, thread)
}

func (db *mysqlDB) GetStakeholderDescendants(thread int64) (map[string]*taps.Personteam, error) {
	return db.getStakeholderAncDes(`
	WITH   RECURSIVE descendants (child, parent) AS
	       (
	       SELECT child
			 ,    parent
	       FROM   threads_parent_child
		   WHERE  parent = %v
	       UNION ALL
	       SELECT t.child
			 ,    t.parent
	       FROM   threads_parent_child t
	         JOIN descendants d
			 ON   t.parent = d.child
		   )
	SELECT DISTINCT stakeholder
	FROM   descendants d
	  JOIN threads_stakeholders t
	  ON   d.child = t.thread
	;`, thread)
}

func (db *mysqlDB) getStakeholderAncDes(query string, thread int64) (map[string]*taps.Personteam, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(query, thread))
	if errQry != nil {
		return nil, fmt.Errorf("Could not query for affected stakeholders: %v", errQry)
	}
	defer qr.Close()
	sks := map[string](*taps.Personteam){}
	for qr.Next() {
		var e string
		errScn := qr.Scan(&e)
		if errScn != nil {
			return nil, fmt.Errorf("Could not scan stakeholder email: %v", errScn)
		}
		pt, errPT := db.GetPersonteam(e)
		if errPT != nil {
			return nil, fmt.Errorf("Could not get personteam from stakeholder email: %v", errPT)
		}
		sks[e] = pt
	}
	return sks, nil
}

func (db *mysqlDB) GetStakeholderOrderBefore(stakeholder, iter string, order int) (int, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_stakeholders
	WHERE  stakeholder = '%v'
	  AND  ord < %v
	  AND  iteration = '%v'
	;`, stakeholder, order, iter))
	if errQry != nil {
		return 0, fmt.Errorf("Could not query for previous thread order: %v", errQry)
	}
	defer qr.Close()
	max := 0
	for qr.Next() {
		errScn := qr.Scan(&max)
		if errScn != nil {
			return 0, fmt.Errorf("Could not scan max value: %v", errScn)
		}
	}
	return max + ((order - max) / 2), nil
}

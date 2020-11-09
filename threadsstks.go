package tapdb

import (
	"fmt"
)

func (db *mysqlDB) NewThreadStkLink(thread int64, stk, domain, iter string, ord int, toplvl bool, cost int) error {
	_, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders
	            (thread,  stk, domain, iter, ord, toplevel, cost)
	VALUES      (    %v, '%v',   '%v', '%v',  %v,       %v,   %v)
	;`, thread, stk, domain, iter, ord, toplvl, cost))
	if errIn != nil {
		return fmt.Errorf("Could not add stakeholder %v to thread %v: %v", stk, thread, errIn)
	}
	return nil
}

func (db *mysqlDB) NewThreadHierLinkForStk(parent, child int64, stk, domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders_hierarchy
				(parent, child,  stk, domain)
	VALUES      (    %v,    %v, '%v',   '%v')
	;`, parent, child, stk, domain))
	return err
}

func (db *mysqlDB) GetOrdBeforeForStk(stk, iter string, ord int) (int, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_stakeholders
	WHERE  stk = %v
	  AND  ord < %v
	  AND  iter = '%v'
	;`, stk, ord, iter))
	if errQry != nil {
		return 0, fmt.Errorf("Could not query for previous thread order: %v", errQry)
	}
	defer qr.Close()
	max := 0
	for qr.Next() {
		errScn := qr.Scan(&max)
		if errScn != nil {
			return 0, nil
		}
	}
	return max, nil
}

func (db *mysqlDB) SetOrdForStk(thread int64, stk string, ord int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    ord = %v
	WHERE  thread = %v
	  AND  stk = %v
	;`, ord, thread, stk))
	return err
}

/*

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

func (db *mysqlDB) SetStakeholderCostTotal(thread int64, stakeholder string, cost int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    costctx = %v
	WHERE  id = %v
	  AND  stakeholder = '%v'
	`, cost, thread, stakeholder))
	return err
}

func (db *mysqlDB) SetStakeholderTopThread(thread int64, stakeholder string, top bool) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    toplevel = %v
	WHERE  id = %v
	  AND  stakeholder = %v
	`, top, thread, stakeholder))
	return err
}
*/

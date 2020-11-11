package tapdb

import (
	"fmt"
)

func (db *mysqlDB) NewThreadStkLink(thread int64, stk, domain, iter string, ord int, toplvl bool, cost int) error {
	_, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders
	            (thread,  stk, domain, iter, ord, toplvl, cost)
	VALUES      (    %v, '%v',   '%v', '%v',  %v,     %v,   %v)
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

func (db *mysqlDB) SetCostForStk(thread int64, stk string, cost int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    cost = %v
	WHERE  thread = %v
	  AND  stk = '%v'
	;`, cost, thread, stk))
	return err
}

func (db *mysqlDB) SetTopForStk(thread int64, stk string, top bool) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    toplvl = %v
	WHERE  id = %v
	  AND  stk = %v
	`, top, thread, stk))
	return err
}
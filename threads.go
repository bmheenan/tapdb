package tapdb

import (
	"fmt"
)

func (db *mysqlDB) NewThread(name, domain, owner, iter, state string, percentile float64, cost int) (int64, error) {
	res, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads
	            (name, domain, owner, iter, state, percentile, costdir, costtot)
	VALUES      ('%v',   '%v',  '%v', '%v',  '%v',         %v,      %v,      %v)
	;`, name, domain, owner, iter, state, percentile, cost, cost))
	if errIn != nil {
		return 0, fmt.Errorf("Could not insert new thread into db: %v", errIn)
	}
	id, errID := res.LastInsertId()
	if errID != nil {
		return 0, fmt.Errorf("Could not get id of inserted thread: %v", errID)
	}
	return id, nil
}

func (db *mysqlDB) NewThreadHierLink(parent, child int64, iter string, ord int, domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_hierarchy
	            (parent, child, domain, iter, ord)
	VALUES      (    %v,    %v,   '%v', '%v',  %v)
	;`, parent, child, domain, iter, ord))
	return err
}

func (db *mysqlDB) NewThreadHierLinkForStk(parent, child int64, stk, domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders_hierarchy
				(parent, child,  stk, domain)
	VALUES      (    %v,    %v, '%v',   '%v')
	;`, parent, child, stk, domain))
	return err
}

func (db *mysqlDB) GetOrdBeforeForParent(parent int64, iter string, ord int) (int, error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_hierarchy
	WHERE  parent = %v
	  AND  ord < %v
	  AND  iteration = '%v'
	;`, parent, ord, iter))
	if errQr != nil {
		return 0, fmt.Errorf("Could not query for previous thread order: %v", errQr)
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

func (db *mysqlDB) SetOrdForParent(thread, parent int64, ord int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_hierarchy
	SET    ord = %v
	WHERE  child = %v
	  AND  parent = %v
	;`, ord, thread, parent))
	return err
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

func (db *mysqlDB) SetCostTot(thread int64, cost int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    costtotal = %v
	WHERE  id = %v
	;`, cost, thread))
	return err
}

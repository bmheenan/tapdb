package tapdb

import (
	"fmt"
	"math"
)

func (db *mysqlDB) NewThreadStkLink(thread int64, stk, domain, iter string, ord int, cost int) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders
	            (thread,  stk, domain, iter, ord, cost)
	VALUES      (    %v, '%v',   '%v', '%v',  %v,   %v)
	;`, thread, stk, domain, iter, ord, cost))
	if err != nil {
		panic(fmt.Sprintf("Could not add stakeholder %v to thread %v: %v", stk, thread, err))
	}
}

/*
func (db *mysqlDB) NewThreadHierLinkForStk(parent, child int64, stk, domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders_hierarchy
				(parent, child,  stk, domain)
	VALUES      (    %v,    %v, '%v',   '%v')
	;`, parent, child, stk, domain))
	return err
}

func (db *mysqlDB) DeleteThreadHierLinkForStk(parent, child int64, stk string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_stakeholders_hierarchy
	WHERE       parent = %v
	  AND       child = %v
	  AND       stk = '%v'
	;`, parent, child, stk))
	return err
}
*/

func (db *mysqlDB) GetOrdBeforeForStk(stk, iter string, ord int) (int, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_stakeholders
	WHERE  stk = '%v'
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

func (db *mysqlDB) GetOrdAfterForStk(stk, iter string, ord int) (int, error) {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT MIN(ord) AS ord
	FROM   threads_stakeholders
	WHERE  stk = '%v'
	  AND  ord > %v
	  AND  iter = '%v'
	;`, stk, ord, iter))
	if err != nil {
		return 0, fmt.Errorf("Could not query for thread order: %v", err)
	}
	defer qr.Close()
	min := 0
	for qr.Next() {
		errScn := qr.Scan(&min)
		if errScn != nil {
			return math.MaxInt32, nil
		}
	}
	return min, nil
}

/*
func (db *mysqlDB) GetChildrenByParentStkLinks(parent int64, stk string) (children []int64, err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT child
	FROM   threads_stakeholders_hierarchy
	WHERE  parent = %v
	  AND  stk = '%v'
	;`, parent, stk))
	if errQr != nil {
		err = fmt.Errorf("Could not query for threads_stakeholders_hierarchy links: %v", errQr)
		return
	}
	defer qr.Close()
	children = []int64{}
	for qr.Next() {
		var c int64
		errScn := qr.Scan(&c)
		if errScn != nil {
			err = fmt.Errorf("Could not scan row: %v", errScn)
			return
		}
		children = append(children, c)
	}
	return
}

func (db *mysqlDB) GetParentsByChildStkLinks(child int64, stk string) (parents []int64, err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT parent
	FROM   threads_stakeholders_hierarchy
	WHERE  child = %v
	  AND  stk = '%v'
	;`, child, stk))
	if errQr != nil {
		err = fmt.Errorf("Could not query for threads_stakeholders_hierarchy links: %v", errQr)
		return
	}
	defer qr.Close()
	parents = []int64{}
	for qr.Next() {
		var p int64
		errScn := qr.Scan(&p)
		if errScn != nil {
			err = fmt.Errorf("Could not scan row: %v", errScn)
			return
		}
		parents = append(parents, p)
	}
	return
}
*/

func (db *mysqlDB) SetOrdForStk(thread int64, stk string, ord int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    ord = %v
	WHERE  thread = %v
	  AND  stk = '%v'
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

func (db *mysqlDB) SetIterForStk(thread int64, stk, iter string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    iter = '%v'
	WHERE  thread = %v
	  AND  stk = '%v'
	;`, iter, thread, stk))
	return err
}

package tapdb

import (
	"fmt"
	"math"
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

func (db *mysqlDB) DeleteThreadHierLink(parent, child int64) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_hierarchy
	WHERE       parent = %v
	  AND       child = %v
	;`, parent, child))
	return err
}

func (db *mysqlDB) GetOrdBeforeForParent(parent int64, iter string, ord int) (int, error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_hierarchy
	WHERE  parent = %v
	  AND  ord < %v
	  AND  iter = '%v'
	;`, parent, ord, iter))
	if errQr != nil {
		return 0, fmt.Errorf("Could not query for thread order: %v", errQr)
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

func (db *mysqlDB) GetOrdAfterForParent(parent int64, iter string, ord int) (int, error) {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT MIN(ord) AS ord
	FROM   threads_hierarchy
	WHERE  parent = %v
	  AND  ord > %v
	  AND  iter = '%v'
	;`, parent, ord, iter))
	if err != nil {
		return 0, fmt.Errorf("Could not query for thread order: %v", err)
	}
	defer qr.Close()
	min := 0
	for qr.Next() {
		err = qr.Scan(&min)
		if err != nil {
			return math.MaxInt32, nil
		}
	}
	return min, nil
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

func (db *mysqlDB) SetCostTot(thread int64, cost int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    costtot = %v
	WHERE  id = %v
	;`, cost, thread))
	return err
}

func (db *mysqlDB) SetIter(thread int64, iter string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    iter = '%v'
	WHERE  id = %v
	;`, iter, thread))
	return err
}

func (db *mysqlDB) SetIterForParent(parent, child int64, iter string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_hierarchy
	SET    iter = '%v'
	WHERE  parent = %v
	  AND  child = %v
	;`, iter, parent, child))
	return err
}

package tapdb

import (
	"fmt"
	"math"

	"github.com/bmheenan/taps"

	"github.com/go-sql-driver/mysql"
)

func (db *mysqlDB) NewThread(name, domain, owner, iter, state string, percentile float64, cost int) int64 {
	res, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads
	            (name, domain, owner, iter, state, percentile, costdir, costtot)
	VALUES      ('%v',   '%v',  '%v', '%v',  '%v',         %v,      %v,      %v)
	;`, name, domain, owner, iter, state, percentile, cost, cost))
	if err != nil {
		panic(fmt.Sprintf("Could not insert new thread into db: %v", err))
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(fmt.Sprintf("Could not get id of inserted thread: %v", err))
	}
	return id
}

func (db *mysqlDB) NewThreadHierLink(parent, child int64, iter string, ord int, domain string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_hierarchy
	            (parent, child, domain, iter, ord)
	VALUES      (    %v,    %v,   '%v', '%v',  %v)
	;`, parent, child, domain, iter, ord))
	if err != nil {
		sqlerr, ok := err.(*mysql.MySQLError)
		if !ok || sqlerr.Number != 1062 { // 1062 = duplicate entry. If they're already linked, don't error
			panic(fmt.Sprintf("Could not insert row linking threads: %v", err))
		}
	}
}

func (db *mysqlDB) DeleteThreadHierLink(parent, child int64) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	DELETE FROM threads_hierarchy
	WHERE       parent = %v
	  AND       child = %v
	;`, parent, child))
	if err != nil {
		panic(fmt.Sprintf("Could not delete thread hierarchy link: %v", err))
	}
}

func (db *mysqlDB) GetOrdBeforeForParent(parent int64, iter string, ord int) int {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_hierarchy
	WHERE  parent = %v
	  AND  ord < %v
	  AND  iter = '%v'
	;`, parent, ord, iter))
	if err != nil {
		panic(fmt.Sprintf("Could not query for thread order: %v", err))
	}
	defer qr.Close()
	max := 0
	for qr.Next() {
		err := qr.Scan(&max)
		if err != nil {
			return 0
		}
	}
	return max
}

func (db *mysqlDB) GetOrdAfterForParent(parent int64, iter string, ord int) int {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT MIN(ord) AS ord
	FROM   threads_hierarchy
	WHERE  parent = %v
	  AND  ord > %v
	  AND  iter = '%v'
	;`, parent, ord, iter))
	if err != nil {
		panic(fmt.Sprintf("Could not query for thread order: %v", err))
	}
	defer qr.Close()
	min := 0
	for qr.Next() {
		err = qr.Scan(&min)
		if err != nil {
			return math.MaxInt32
		}
	}
	return min
}

func (db *mysqlDB) SetOrdForParent(thread, parent int64, ord int) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_hierarchy
	SET    ord = %v
	WHERE  child = %v
	  AND  parent = %v
	;`, ord, thread, parent))
	if err != nil {
		panic(fmt.Sprintf("Could not set ord: %v", err))
	}
}

func (db *mysqlDB) SetCostTot(thread int64, cost int) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    costtot = %v
	WHERE  id = %v
	;`, cost, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set cost: %v", err))
	}
}

func (db *mysqlDB) SetIter(thread int64, iter string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    iter = '%v'
	WHERE  id = %v
	;`, iter, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set iter: %v", err))
	}
}

func (db *mysqlDB) SetIterForParent(parent, child int64, iter string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_hierarchy
	SET    iter = '%v'
	WHERE  parent = %v
	  AND  child = %v
	;`, iter, parent, child))
	if err != nil {
		panic(fmt.Sprintf("Could not set iter for parent: %v", err))
	}
}

func (db *mysqlDB) SetName(thread int64, name string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    name = '%v'
	WHERE  id = %v
	;`, name, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set name: %v", err))
	}
}

func (db *mysqlDB) SetDesc(thread int64, desc string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    description = '%v'
	WHERE  id = %v
	;`, desc, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set description: %v", err))
	}
}

func (db *mysqlDB) SetCostDir(thread int64, cost int) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    costdir = %v
	WHERE  id = %v
	;`, cost, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set cost: %v", err))
	}
}

func (db *mysqlDB) SetState(thread int64, state taps.State) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    state = '%v'
	WHERE  id = %v
	;`, state, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set state: %v", err))
	}
}

func (db *mysqlDB) SetOwner(thread int64, owner string) {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    owner = '%v'
	WHERE  id = %v
	;`, owner, thread))
	if err != nil {
		panic(fmt.Sprintf("Could not set state: %v", err))
	}
}

package tapdb

import (
	"fmt"
)

func (db *mysqlDB) GetItersForStk(stk string) (iters []string, err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   DISTINCT iter
	FROM     threads_stakeholders
	WHERE    stk = '%v'
	ORDER BY iter
	;`, stk))
	if errQr != nil {
		err = fmt.Errorf("Could not query for iters for %v: %v", stk, errQr)
		return
	}
	defer qr.Close()
	iters = []string{}
	for qr.Next() {
		var iter string
		errScn := qr.Scan(&iter)
		if errScn != nil {
			err = fmt.Errorf("Could not scan iteration: %v", errScn)
			return
		}
		iters = append(iters, iter)
	}
	return
}

func (db *mysqlDB) GetItersForParent(parent int64) (iters []string, err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   DISTINCT iter
	FROM     threads_hierarchy
	WHERE    parent = %v
	ORDER BY iter
	;`, parent))
	if errQr != nil {
		err = fmt.Errorf("Could not query for iters for %v: %v", parent, errQr)
		return
	}
	defer qr.Close()
	iters = []string{}
	for qr.Next() {
		var iter string
		errScn := qr.Scan(&iter)
		if errScn != nil {
			err = fmt.Errorf("Could not scan iteration: %v", errScn)
			return
		}
		iters = append(iters, iter)
	}
	return
}

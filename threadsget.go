package tapdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmheenan/taps"
)

func (db *mysqlDB) GetThread(thread int64) (*taps.Thread, error) {
	thQr, errTh := db.conn.Query(fmt.Sprintf(`
	SELECT id
	  ,    name
	  ,    description
	  ,    state
	  ,    costdir
	  ,    costtot
	  ,    owner
	  ,    iter
	  ,    percentile
	FROM   threads
	WHERE  id = %v
	;`, thread))
	if errTh != nil {
		return nil, fmt.Errorf("Could not query data for thread %v: %v", thread, errTh)
	}
	defer thQr.Close()
	if thQr.Next() {
		var oEmail string
		var desc sql.NullString
		th := &taps.Thread{}
		errScn := thQr.Scan(
			&th.ID,
			&th.Name,
			&desc,
			&th.State,
			&th.CostDir,
			&th.CostTot,
			&oEmail,
			&th.Iter,
			&th.Percentile,
		)
		if errScn != nil {
			return nil, fmt.Errorf("Could not scan thread %v: %v", thread, errScn)
		}
		if desc.Valid {
			th.Desc = desc.String
		}
		o, errO := db.GetStk(oEmail)
		if errO != nil {
			return nil, fmt.Errorf("Could not get stakeholder %v: %v", oEmail, errO)
		}
		th.Owner = *o
		th.Stks = map[string](struct {
			Iter string
			Ord  int
			Cost int
		}){}
		stksQr, errStks := db.conn.Query(fmt.Sprintf(`
		SELECT stk
		  ,    iter
		  ,    ord
		  ,    cost
		FROM   threads_stakeholders
		WHERE  thread = %v
		;`, th.ID))
		if errStks != nil {
			return nil, fmt.Errorf("Could not query stakeholders of %v: %v", th.ID, errStks)
		}
		defer stksQr.Close()
		for stksQr.Next() {
			var e string
			stk := struct {
				Iter string
				Ord  int
				Cost int
			}{}
			errScn := stksQr.Scan(&e, &stk.Iter, &stk.Ord, &stk.Cost)
			if errScn != nil {
				return nil, fmt.Errorf("Could not scan stakeholder: %v", errScn)
			}
			th.Stks[e] = stk
		}
		th.Parents = map[int64](struct {
			Iter string
			Ord  int
		}){}
		parQr, errP := db.conn.Query(fmt.Sprintf(`
		SELECT parent
		  ,    iter
		  ,    ord
		FROM   threads_hierarchy
		WHERE  child = %v
		;`, th.ID))
		if errP != nil {
			return nil, fmt.Errorf("Could not query parents of %v: %v", th.ID, errP)
		}
		defer parQr.Close()
		for parQr.Next() {
			var i int64
			p := struct {
				Iter string
				Ord  int
			}{}
			errScn := parQr.Scan(&i, &p.Iter, &p.Ord)
			if errScn != nil {
				return nil, fmt.Errorf("Could not scan parent: %v", errScn)
			}
			th.Parents[i] = p
		}
		return th, nil
	}
	return nil, fmt.Errorf("No thread found with id %v: %w", thread, ErrNotFound)
}

func (db *mysqlDB) GetThreadsByStkIter(stk, iter string) ([](*taps.Thread), error) {
	return db.getThreadsByStkPaIter(0, stk, iter)
}

func (db *mysqlDB) GetThreadsByParentIter(parent int64, iter string) ([](*taps.Thread), error) {
	return db.getThreadsByStkPaIter(parent, "", iter)
}

func (db *mysqlDB) getThreadsByStkPaIter(parent int64, stk, iter string) ([](*taps.Thread), error) {
	var (
		qr    *sql.Rows
		errQr error
	)
	if stk != "" {
		qr, errQr = db.conn.Query(fmt.Sprintf(`
		SELECT   thread
		FROM     threads_stakeholders
		WHERE    stk = '%v'
		  AND    iter = '%v'
		ORDER BY ord
		;`, stk, iter))
	} else if parent != 0 {
		qr, errQr = db.conn.Query(fmt.Sprintf(`
		SELECT   child
		FROM     threads_hierarchy
		WHERE    parent = %v
		  AND    iter = '%v'
		ORDER BY ord
		;`, parent, iter))
	} else {
		return nil, errors.New("stk must != '' or parent must != 0")
	}
	if errQr != nil {
		return nil, fmt.Errorf("Could not query for threads: %v", errQr)
	}
	defer qr.Close()
	ids := []int64{}
	for qr.Next() {
		var id int64
		errScn := qr.Scan(&id)
		if errScn != nil {
			return nil, fmt.Errorf("Could not scan id: %v", errScn)
		}
		ids = append(ids, id)
	}
	ths := [](*taps.Thread){}
	for _, id := range ids {
		th, errTh := db.GetThread(id)
		if errTh != nil {
			return nil, fmt.Errorf("Could not get thread %v: %v", id, errTh)
		}
		ths = append(ths, th)
	}
	return ths, nil
}

func (db *mysqlDB) GetThreadDes(thread int64) (map[int64](*taps.Thread), error) {
	thTop, errTop := db.GetThread(thread)
	if errTop != nil {
		return nil, fmt.Errorf("Could not get top thread: %w", errTop)
	}
	ths := map[int64](*taps.Thread){
		thTop.ID: thTop,
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	WITH   RECURSIVE des (child, parent) AS
	       (
	       SELECT child
	         ,    parent
	       FROM   threads_hierarchy
	       WHERE  parent = %v
	       UNION ALL
	       SELECT t.child
	         ,    t.parent
	       FROM   threads_hierarchy t
	       JOIN   des d
	         ON   t.parent = d.child
		   )
	SELECT DISTINCT child
	FROM   des
	;`, thread))
	if errQry != nil {
		return nil, fmt.Errorf("Could not query for descendant threads: %v", errQry)
	}
	defer qr.Close()
	for qr.Next() {
		var i int64
		qr.Scan(&i)
		th, errTh := db.GetThread(i)
		if errTh != nil {
			return nil, fmt.Errorf("Could not get descendant thread: %v", errTh)
		}
		ths[th.ID] = th
	}
	return ths, nil
}

func (db *mysqlDB) GetThreadAns(thread int64) (map[int64](*taps.Thread), error) {
	thBtm, errBtm := db.GetThread(thread)
	if errBtm != nil {
		return nil, fmt.Errorf("Could not get bottom thread: %w", errBtm)
	}
	ths := map[int64](*taps.Thread){
		thBtm.ID: thBtm,
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	WITH   RECURSIVE ans (child, parent) AS
	       (
	       SELECT child
	         ,    parent
	       FROM   threads_hierarchy
	       WHERE  child = %v
	       UNION ALL
	       SELECT t.child
	         ,    t.parent
	       FROM   threads_hierarchy t
	       JOIN   ans a
	         ON   t.child = a.parent
		   )
	SELECT DISTINCT parent
	FROM   ans
	;`, thread))
	if errQry != nil {
		return nil, fmt.Errorf("Could not query for ancestor threads: %v", errQry)
	}
	defer qr.Close()
	for qr.Next() {
		var i int64
		qr.Scan(&i)
		th, errTh := db.GetThread(i)
		if errTh != nil {
			return nil, fmt.Errorf("Could not get ancestor thread: %v", errTh)
		}
		ths[th.ID] = th
	}
	return ths, nil
}

func (db *mysqlDB) GetThreadChildrenByStkIter(threads []int64, stk, iter string) (map[int64](*taps.Thread), error) {
	return db.getThChPaByStkIter(threads, stk, iter, "children")
}

func (db *mysqlDB) GetThreadParentsByStkIter(threads []int64, stk, iter string) (map[int64](*taps.Thread), error) {
	return db.getThChPaByStkIter(threads, stk, iter, "parents")
}

func (db *mysqlDB) getThChPaByStkIter(threads []int64, stk, iter, dir string) (map[int64](*taps.Thread), error) {
	ret := map[int64](*taps.Thread){}
	for _, id := range threads {
		th, errTh := db.GetThread(id)
		if errTh != nil {
			return map[int64](*taps.Thread){}, fmt.Errorf("Could not get thread from id %v: %v", id, errTh)
		}
		if _, ok := th.Stks[stk]; ok && th.Iter == iter {
			ret[id] = th
		} else {
			var (
				qr     *sql.Rows
				errQry error
			)
			if dir == "children" {
				qr, errQry = db.conn.Query(fmt.Sprintf(`
				SELECT child
				FROM   threads_hierarchy
				WHERE  parent = %v
				;`, id))
			} else if dir == "parents" {
				qr, errQry = db.conn.Query(fmt.Sprintf(`
				SELECT parent
				FROM   threads_hierarchy
				WHERE  child = %v
				;`, id))
			}
			if errQry != nil {
				return map[int64](*taps.Thread){}, fmt.Errorf(
					"Could not query %v of id %v: %v",
					dir,
					id,
					errQry,
				)
			}
			defer qr.Close()
			cids := []int64{}
			for qr.Next() {
				var cid int64
				errScn := qr.Scan(&cid)
				if errScn != nil {
					return map[int64](*taps.Thread){}, fmt.Errorf("Could not scan query result: %v", errScn)
				}
				cids = append(cids, cid)
			}
			lc, errLC := db.getThChPaByStkIter(cids, stk, iter, dir)
			if errLC != nil {
				return map[int64](*taps.Thread){}, fmt.Errorf(
					"Could not get stakeholder %v of %v: %v",
					dir,
					cids,
					errLC,
				)
			}
			for i, c := range lc {
				ret[i] = c
			}
		}
	}
	return ret, nil
}

func (db *mysqlDB) GetThreadrowsByStkIter(stk, iter string) ([](*taps.Threadrow), error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT      t.id
	  ,         t.name
	  ,         t.state
	  ,         s.cost
	  ,         t.owner
	  ,         t.iter
	  ,         s.ord
	FROM        threads t
	  JOIN      threads_stakeholders s
	  ON        t.id = s.thread
	  LEFT JOIN threads_stakeholders_hierarchy h
	  ON        t.id = h.child
	WHERE       h.child IS NULL
	  AND       s.stk = '%v'
	  AND       s.iter = '%v'
	;`, stk, iter))
	if errQr != nil {
		return nil, fmt.Errorf("Could not query for top level threads: %v", errQr)
	}
	defer qr.Close()
	ths := [](*taps.Threadrow){}
	for qr.Next() {
		th := &taps.Threadrow{}
		var oEmail string
		errScn := qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter, &th.Ord)
		if errScn != nil {
			return nil, fmt.Errorf("Could not scan top level thread: %v", errScn)
		}
		tOwner, errO := db.GetStk(oEmail)
		if errO != nil {
			return nil, fmt.Errorf("Could not get stakeholder from email %v: %v", oEmail, errO)
		}
		th.Owner = *tOwner
		errDes := db.fillThreadrowDesByStkIter(th, stk, iter)
		if errDes != nil {
			return nil, fmt.Errorf("Could not fill decendants of %v: %v", th.Name, errDes)
		}
		ths = append(ths, th)
	}
	return ths, nil
}

func (db *mysqlDB) fillThreadrowDesByStkIter(parent *taps.Threadrow, stk, iter string) error {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   h.child
	  ,      t.name
	  ,      t.state
	  ,      s.cost
	  ,      t.owner
	  ,      t.iter
	  ,      s.ord
	FROM     threads_stakeholders_hierarchy h
	  JOIN   threads t
	  ON     h.child = t.id
	  JOIN   threads_stakeholders s
	  ON     h.child = s.thread
	    AND  h.stk = s.stk 
	WHERE    h.parent = %v
	  AND    h.stk = '%v'
	  AND    s.iter = '%v'
	ORDER BY s.ord
	;`, parent.ID, stk, iter))
	if errQr != nil {
		return fmt.Errorf("Could not query for children: %v", errQr)
	}
	defer qr.Close()
	for qr.Next() {
		th := taps.Threadrow{}
		var oEmail string
		errScn := qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter, &th.Ord)
		if errScn != nil {
			return fmt.Errorf("Could not scan thread: %v", errScn)
		}
		thO, errO := db.GetStk(oEmail)
		if errO != nil {
			return fmt.Errorf("Could not get stakeholder from email %v: %v", oEmail, errO)
		}
		th.Owner = *thO
		parent.Children = append(parent.Children, th)
	}
	return nil
}

func (db *mysqlDB) GetThreadrowsByParentIter(parent int64, iter string) ([](*taps.Threadrow), error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	WITH     q AS
	         (
		     SELECT child
			   ,    iter
			   ,    ord
	         FROM   threads_hierarchy
	         WHERE  parent = %v
	           AND  iter = '%v'
	         )
	SELECT   t.id
	  ,      t.name
	  ,      t.state
	  ,      t.costtot
	  ,      t.owner
	  ,      t.iter
	  ,      q.ord
	FROM     threads t
	  JOIN   q
	  ON     t.id = q.child
	ORDER BY q.ord
	;`, parent, iter))
	if errQr != nil {
		return nil, fmt.Errorf("Could not query for child threads of %v: %v", parent, errQr)
	}
	defer qr.Close()
	ths := [](*taps.Threadrow){}
	for qr.Next() {
		th := &taps.Threadrow{}
		var oEmail string
		errScn := qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter, &th.Ord)
		if errScn != nil {
			return nil, fmt.Errorf("Could not scan child thread: %v", errScn)
		}
		tOwner, errO := db.GetStk(oEmail)
		if errO != nil {
			return nil, fmt.Errorf("Could not get stakeholder from email %v: %v", oEmail, errO)
		}
		th.Owner = *tOwner
		errDes := db.fillThreadrowDes(th)
		if errDes != nil {
			return nil, fmt.Errorf("Could not fill decendants of %v: %v", th.Name, errDes)
		}
		ths = append(ths, th)
	}
	return ths, nil
}

func (db *mysqlDB) fillThreadrowDes(parent *taps.Threadrow) error {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   h.child
	  ,      t.name
	  ,      t.state
	  ,      t.costtot
	  ,      t.owner
	  ,      t.iter
	  ,      h.ord
	FROM     threads_hierarchy h
	  JOIN   threads t
	  ON     h.child = t.id
	WHERE    h.parent = %v
	ORDER BY h.ord
	;`, parent.ID))
	if errQr != nil {
		return fmt.Errorf("Could not query for children: %v", errQr)
	}
	defer qr.Close()
	for qr.Next() {
		th := taps.Threadrow{}
		var oEmail string
		errScn := qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter, &th.Ord)
		if errScn != nil {
			return fmt.Errorf("Could not scan thread: %v", errScn)
		}
		thO, errO := db.GetStk(oEmail)
		if errO != nil {
			return fmt.Errorf("Could not get stakeholder from email %v: %v", oEmail, errO)
		}
		th.Owner = *thO
		parent.Children = append(parent.Children, th)
	}
	return nil
}

package tapdb

import (
	"database/sql"
	"fmt"
	"sort"

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

func (db *mysqlDB) GetThreadsByStkIter(stk, iter string) []*taps.Thread {
	return db.getThreadsByStkPaIter(0, stk, iter)
}

func (db *mysqlDB) GetThreadsByParentIter(parent int64, iter string) []*taps.Thread {
	return db.getThreadsByStkPaIter(parent, "", iter)
}

func (db *mysqlDB) getThreadsByStkPaIter(parent int64, stk, iter string) []*taps.Thread {
	var (
		qr  *sql.Rows
		err error
	)
	if stk != "" {
		qr, err = db.conn.Query(fmt.Sprintf(`
		SELECT   thread
		FROM     threads_stakeholders
		WHERE    stk = '%v'
		  AND    iter = '%v'
		ORDER BY ord
		;`, stk, iter))
	} else if parent != 0 {
		qr, err = db.conn.Query(fmt.Sprintf(`
		SELECT   child
		FROM     threads_hierarchy
		WHERE    parent = %v
		  AND    iter = '%v'
		ORDER BY ord
		;`, parent, iter))
	} else {
		panic("stk must != '' or parent must != 0")
	}
	if err != nil {
		panic(fmt.Sprintf("Could not query for threads: %v", err))
	}
	defer qr.Close()
	ids := []int64{}
	for qr.Next() {
		var id int64
		err = qr.Scan(&id)
		if err != nil {
			panic(fmt.Sprintf("Could not scan id: %v", err))
		}
		ids = append(ids, id)
	}
	ths := [](*taps.Thread){}
	for _, id := range ids {
		th, err := db.GetThread(id)
		if err != nil {
			panic(fmt.Sprintf("Could not get thread %v: %v", id, err))
		}
		ths = append(ths, th)
	}
	return ths
}

func (db *mysqlDB) GetThreadDes(thread int64) map[int64]*taps.Thread {
	thTop, err := db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get top thread: %v", err))
	}
	ths := map[int64](*taps.Thread){
		thTop.ID: thTop,
	}
	qr, err := db.conn.Query(fmt.Sprintf(`
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
	if err != nil {
		panic(fmt.Sprintf("Could not query for descendant threads: %v", err))
	}
	defer qr.Close()
	for qr.Next() {
		var i int64
		qr.Scan(&i)
		th, err := db.GetThread(i)
		if err != nil {
			panic(fmt.Sprintf("Could not get descendant thread: %v", err))
		}
		ths[th.ID] = th
	}
	return ths
}

func (db *mysqlDB) GetThreadAns(thread int64) map[int64]*taps.Thread {
	thBtm, err := db.GetThread(thread)
	if err != nil {
		panic(fmt.Sprintf("Could not get bottom thread: %v", err))
	}
	ths := map[int64](*taps.Thread){
		thBtm.ID: thBtm,
	}
	qr, err := db.conn.Query(fmt.Sprintf(`
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
	if err != nil {
		panic(fmt.Sprintf("Could not query for ancestor threads: %v", err))
	}
	defer qr.Close()
	for qr.Next() {
		var i int64
		qr.Scan(&i)
		th, err := db.GetThread(i)
		if err != nil {
			panic(fmt.Sprintf("Could not get ancestor thread: %v", err))
		}
		ths[th.ID] = th
	}
	return ths
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

func (db *mysqlDB) GetThreadrowsByStkIter(stk, iter string) (ths []taps.Threadrow) {
	qr, err := db.conn.Query(fmt.Sprintf(`
	WITH        RECURSIVE des (parent, child) AS
	            (
				SELECT h.parent
				  ,    h.child
				FROM   threads_hierarchy h
				JOIN   threads_stakeholders s
				  ON   h.parent = s.thread
				WHERE  s.stk = '%v'
				  AND  s.iter = '%v'
				UNION ALL
				SELECT h.parent
				  ,    h.child
				FROM   threads_hierarchy h
				  JOIN des d
				  ON   h.parent = d.child
				)
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
	  LEFT JOIN des d
	  ON        t.id = d.child
	WHERE       d.child IS NULL
	  AND       s.stk = '%v'
	  AND       s.iter = '%v'
	ORDER BY    s.ord
	;`, stk, iter, stk, iter))
	if err != nil {
		panic(fmt.Sprintf("Could not query for top level threadrows: %v", err))
	}
	defer qr.Close()
	ths = []taps.Threadrow{}
	for qr.Next() {
		th := taps.Threadrow{}
		oe := ""
		err = qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oe, &th.Iter, &th.Ord)
		if err != nil {
			panic(fmt.Sprintf("Could not scan threadrows: %v", err))
		}
		o, err := db.GetStk(oe)
		if err != nil {
			panic(fmt.Sprintf("Could not get stakeholder for owner of thread: %v", err))
		}
		th.Owner = *o
		db.fillThreadrowDesByStkIter(th.ID, &th.Children, stk, iter, &map[int64]string{})
		sort.Slice(th.Children, func(i, j int) bool {
			return th.Children[i].Ord < th.Children[j].Ord
		})
		ths = append(ths, th)
	}
	return ths
}

func (db *mysqlDB) fillThreadrowDesByStkIter(paID int64, children *[]taps.Threadrow, stk, iter string, added *map[int64]string) {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT      t.id
	  ,         t.name
	  ,         t.state
	  ,         t.owner
	  ,         t.iter
	  ,         s.stk
	  ,         s.iter
	  ,         s.cost
	  ,         s.ord
	FROM        threads_hierarchy h
	  JOIN      threads t
	  ON        h.child = t.id
	  LEFT JOIN threads_stakeholders s
	  ON        t.id = s.thread
	WHERE       h.parent = %v
	ORDER BY    s.ord
	;`, paID))
	if err != nil {
		panic(fmt.Sprintf("Could not query for thread children: %v", err))
	}
	defer qr.Close()
	skippedThs := map[int64]string{}
	for qr.Next() {
		th := taps.Threadrow{}
		var (
			oEmail string
			sEmail,
			sIter sql.NullString
			sCost,
			sOrd sql.NullInt32
		)
		err := qr.Scan(&th.ID, &th.Name, &th.State, &oEmail, &th.Iter, &sEmail, &sIter, &sCost, &sOrd)
		if err != nil {
			panic(fmt.Sprintf("Could not scan thread: %v", err))
		}
		if sEmail.String == stk && sIter.String == iter {
			owner, err := db.GetStk(oEmail)
			if err != nil {
				panic(fmt.Sprintf("Could not get stakeholder from email %v: %v", oEmail, err))
			}
			th.Owner = *owner
			th.Cost = int(sCost.Int32)
			th.Ord = int(sOrd.Int32)
			db.fillThreadrowDesByStkIter(th.ID, &th.Children, stk, iter, added)
			sort.Slice(th.Children, func(i, j int) bool {
				return th.Children[i].Ord < th.Children[j].Ord
			})
			*children = append(*children, th)
			(*added)[th.ID] = th.Name
		} else {
			skippedThs[th.ID] = th.Name
		}
		for id, th := range skippedThs {
			if _, ok := (*added)[id]; !ok {
				(*added)[id] = th
				db.fillThreadrowDesByStkIter(id, children, stk, iter, added)
			}
		}
	}
}

func (db *mysqlDB) GetThreadrowsByParentIter(parent int64, iter string) []taps.Threadrow {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT   t.id
	  ,      t.name
	  ,      t.state
	  ,      t.costtot
	  ,      t.owner
	  ,      t.iter
	  ,      h.ord
	FROM     threads t
	  JOIN   threads_hierarchy h
	  ON     t.id = h.child
	WHERE    h.parent = %v
	  AND    h.iter = '%v'
	ORDER BY h.ord
	;`, parent, iter))
	if err != nil {
		panic(fmt.Sprintf("Could not query for child threads of %v: %v", parent, err))
	}
	defer qr.Close()
	ths := []taps.Threadrow{}
	for qr.Next() {
		th := taps.Threadrow{}
		var oEmail string
		err := qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter, &th.Ord)
		if err != nil {
			panic(fmt.Sprintf("Could not scan child thread: %v", err))
		}
		tOwner, err := db.GetStk(oEmail)
		if err != nil {
			panic(fmt.Sprintf("Could not get stakeholder from email %v: %v", oEmail, err))
		}
		th.Owner = *tOwner
		err = db.fillThreadrowDes(&th)
		if err != nil {
			panic(fmt.Sprintf("Could not fill decendants of %v: %v", th.Name, err))
		}
		ths = append(ths, th)
	}
	return ths
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
	ORDER BY t.iter
	  ,      h.ord
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
		errDes := db.fillThreadrowDes(&th)
		if errDes != nil {
			return fmt.Errorf("Could not fill decendants of %v: %v", th.Name, errDes)
		}
		parent.Children = append(parent.Children, th)
	}
	return nil
}

func (db *mysqlDB) GetThreadParentsForAnc(child, anc int64) (parents []*taps.Thread, err error) {
	parents = []*taps.Thread{}
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT h.parent
	FROM   threads_hierarchy h
	WHERE  h.child = '%v'
	;`, child))
	if err != nil {
		err = fmt.Errorf("Could not get parents of %v: %v", child, err)
		return
	}
	defer qr.Close()
	for qr.Next() {
		var pid int64
		err = qr.Scan(&pid)
		if err != nil {
			err = fmt.Errorf("Could not scan parent id: %v", err)
			return
		}
		ans := db.GetThreadAns(pid)
		if th, ok := ans[anc]; ok {
			parents = append(parents, th)
		}
	}
	return
}

func (db *mysqlDB) GetThreadrowsByChild(child int64) []taps.Threadrow {
	qr, err := db.conn.Query(fmt.Sprintf(`
	SELECT h.parent
	  ,    t.name
	  ,    t.state
	  ,    t.costtot
	  ,    t.owner
	  ,    t.iter
	FROM   threads_hierarchy h
	  JOIN threads t
	  ON   h.parent = t.id
	WHERE  h.child = '%v'
	;`, child))
	if err != nil {
		panic(fmt.Sprintf("Could not get parents of %v: %v", child, err))
	}
	defer qr.Close()
	ths := []taps.Threadrow{}
	for qr.Next() {
		var (
			th     taps.Threadrow
			oEmail string
		)
		err = qr.Scan(&th.ID, &th.Name, &th.State, &th.Cost, &oEmail, &th.Iter)
		if err != nil {
			panic(fmt.Sprintf("Could not scan thread: %v", err))
		}
		thOwner, err := db.GetStk(oEmail)
		if err != nil {
			panic(fmt.Sprintf("Could not get owner from email: %v", err))
		}
		th.Owner = *thOwner
		ths = append(ths, th)
	}
	return ths
}

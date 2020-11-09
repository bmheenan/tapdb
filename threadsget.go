package tapdb

import (
	"database/sql"
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
			Iter   string
			Ord    int
			Cost   int
			Toplvl bool
		}){}
		stksQr, errStks := db.conn.Query(fmt.Sprintf(`
		SELECT stk
		  ,    iter
		  ,    ord
		  ,    cost
		  ,    toplvl
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
				Iter   string
				Ord    int
				Cost   int
				Toplvl bool
			}{}
			errScn := stksQr.Scan(&e, &stk.Iter, &stk.Ord, &stk.Cost, &stk.Toplvl)
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

/*
// GetThreadDescendants returns a threadrel map containing all threads that are descendands of the provided `id`
// (including itself). `StakeholderMatch` and `Order` will not be filled in
func (db *mysqlDB) GetThreadDescendants(id int64, stakeholder string) (map[int64](*taps.Threadrel), error) {
	// TODO This function and GetThreadAncestors can be cleaned up to remove scoping to a single stakeholder and having
	// Threadrel have a map[stakeholder email](iteration, order)
	thTop, errTop := db.GetThreadrel(id, stakeholder)
	if errTop != nil {
		return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not get the root threadrel: %w", errTop)
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	WITH        RECURSIVE descendants (child, parent) AS
	            (
	            SELECT child
	              ,    parent
	            FROM   threads_parent_child
	            WHERE  parent = %v
	            UNION ALL
	            SELECT t.child
	              ,    t.parent
	            FROM   threads_parent_child t
	            JOIN   descendants d
	              ON   t.parent = d.child
		        )
	  ,         filteredstakeholders (thread, stakeholder, ord) AS
	            (
		        SELECT thread
			      ,    stakeholder
			      ,    ord
		        FROM   threads_stakeholders
		        WHERE  stakeholder = '%v'
		        )
	SELECT      t.id
	  ,         t.state
	  ,         t.costdirect
	  ,         t.owner
	  ,         t.iteration
	  ,         t.percentile
	  ,         CASE WHEN s.stakeholder = '%v'
					 THEN true
					 ELSE false
					 END AS stakeholdermatch
	  ,         CASE WHEN s.ord IS NULL
					 THEN 0
					 ELSE s.ord
					 END AS ord
	FROM        descendants d
	  JOIN      threads t
	  ON        t.id = d.child
	  LEFT JOIN filteredstakeholders s
	  ON        t.id = s.thread
	;`, id, stakeholder, stakeholder))
	if errQry != nil {
		return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not query descendants: %v", errQry)
	}
	defer qr.Close()
	ths := map[int64](*taps.Threadrel){
		id: thTop,
	}
	for qr.Next() {
		th := &taps.Threadrel{}
		qr.Scan(
			&th.ID,
			&th.State,
			&th.CostDirect,
			&th.Owner,
			&th.Iteration,
			&th.Percentile,
			&th.StakeholderMatch,
			&th.Order,
		)
		ths[th.ID] = th
	}
	return ths, nil
}

func (db *mysqlDB) GetThreadAncestors(id int64, stakeholder string) (map[int64](*taps.Threadrel), error) {
	thRoot, errRoot := db.GetThreadrel(id, stakeholder)
	if errRoot != nil {
		return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not get the root threadrel: %w", errRoot)
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	WITH        RECURSIVE ancestors (child, parent) AS
	            (
	            SELECT child
	              ,    parent
	            FROM   threads_parent_child
	            WHERE  child = %v
	            UNION ALL
	            SELECT t.child
	              ,    t.parent
	            FROM   threads_parent_child t
	            JOIN   ancestors a
	              ON   t.child = a.parent
		        )
	  ,         filteredstakeholders (thread, stakeholder, ord) AS
	            (
		        SELECT thread
		          ,    stakeholder
		          ,    ord
		        FROM   threads_stakeholders
		        WHERE  stakeholder = '%v'
		        )
	SELECT      t.id
	  ,         t.state
	  ,         t.costdirect
	  ,         t.owner
	  ,         t.iteration
	  ,         t.percentile
	  ,         CASE WHEN s.stakeholder = '%v'
					 THEN true
					 ELSE false
					 END AS stakeholdermatch
	  ,         CASE WHEN s.ord IS NULL
					 THEN 0
					 ELSE s.ord
					 END AS ord
	FROM        ancestors a
	  JOIN      threads t
	  ON        t.id = a.parent
	  LEFT JOIN filteredstakeholders s
	  ON        t.id = s.thread
	;`, id, stakeholder, stakeholder))
	if errQry != nil {
		return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not query ancestors: %v", errQry)
	}
	defer qr.Close()
	ths := map[int64](*taps.Threadrel){
		id: thRoot,
	}
	for qr.Next() {
		th := &taps.Threadrel{}
		errScn := qr.Scan(
			&th.ID,
			&th.State,
			&th.CostDirect,
			&th.Owner,
			&th.Iteration,
			&th.Percentile,
			&th.StakeholderMatch,
			&th.Order,
		)
		if errScn != nil {
			return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not scan from query results: %v", errScn)
		}
		ths[th.ID] = th
	}
	return ths, nil
}

// GetChildThreadsSkIter returns the topmost thread(s) within `threads` and all their descendants, who have
// `stakeholder` as a stakeholder and are in `iteration`
func (db *mysqlDB) GetChildThreadsSkIter(
	threads []int64,
	stakeholder,
	iteration string,
) (map[int64](*taps.Threadrel), error) {
	return db.getChPaThreadsSkIter(threads, stakeholder, iteration, "children")
}

// GetParentThreadsSkIter returns the bottommost thread(s) within `threads` and all their ancestors, who have
// `stakeholder` as a stakeholder and are in `iteration`
func (db *mysqlDB) GetParentThreadsSkIter(
	threads []int64,
	stakeholder,
	iteration string,
) (map[int64](*taps.Threadrel), error) {
	return db.getChPaThreadsSkIter(threads, stakeholder, iteration, "parents")
}

func (db *mysqlDB) getChPaThreadsSkIter(
	threads []int64,
	stakeholder,
	iteration,
	direction string,
) (map[int64](*taps.Threadrel), error) {
	ret := map[int64](*taps.Threadrel){}
	for _, id := range threads {
		th, errTh := db.GetThreadrel(id, stakeholder)
		if errTh != nil {
			return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not get thread from id %v: %v", id, errTh)
		}
		if th.StakeholderMatch && th.Iteration == iteration {
			ret[id] = th
		} else {
			var (
				qr     *sql.Rows
				errQry error
			)
			if direction == "children" {
				qr, errQry = db.conn.Query(fmt.Sprintf(`
				SELECT child
				FROM   threads_parent_child
				WHERE  parent = %v
				;`, id))
			} else if direction == "parents" {
				qr, errQry = db.conn.Query(fmt.Sprintf(`
				SELECT parent
				FROM   threads_parent_child
				WHERE  child = %v
				;`, id))
			}
			if errQry != nil {
				return map[int64](*taps.Threadrel){}, fmt.Errorf(
					"Could not query %v of id %v: %v",
					direction,
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
					return map[int64](*taps.Threadrel){}, fmt.Errorf("Could not scan query result: %v", errScn)
				}
				cids = append(cids, cid)
			}
			lc, errLC := db.getChPaThreadsSkIter(cids, stakeholder, iteration, direction)
			if errLC != nil {
				return map[int64](*taps.Threadrel){}, fmt.Errorf(
					"Could not get stakeholder %v of %v: %v",
					direction,
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

func (db *mysqlDB) GetThreadrelsByParentIter(parent int64, iter string) ([](*taps.Threadrel), error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   t.id
	  ,      t.state
	  ,      t.owner
	  ,      t.percentile
	  ,      t.iteration
	  ,      t.costdirect
	  ,      pc.ord
	FROM     threads t
	  JOIN   threads_parent_child pc
	  ON     pc.child = t.id
	WHERE    pc.parent = %v
	  AND    pc.iteration = '%v'
	ORDER BY pc.ord
	;`, parent, iter))
	if errQr != nil {
		return [](*taps.Threadrel){}, fmt.Errorf("Could not query for threads: %v", errQr)
	}
	defer qr.Close()
	ths := [](*taps.Threadrel){}
	for qr.Next() {
		th := &taps.Threadrel{}
		errScn := qr.Scan(&th.ID, &th.State, &th.Owner, &th.Percentile, &th.Iteration, &th.CostDirect, &th.Order)
		if errScn != nil {
			return [](*taps.Threadrel){}, fmt.Errorf("Could not scan thread: %v", errScn)
		}
		ths = append(ths, th)
	}
	return ths, nil
}

func (db *mysqlDB) GetThreadrelsByStakeholderIter(stakeholder, iter string) ([](*taps.Threadrel), error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   t.id
	  ,      t.state
	  ,      t.owner
	  ,      t.percentile
	  ,      t.iteration
	  ,      t.costdirect
	  ,      s.ord
	FROM     threads t
	  JOIN   threads_stakeholders s
	  ON     s.thread = t.id
	WHERE    s.stakeholder = '%v'
	  AND    s.iteration = '%v'
	ORDER BY s.ord
	;`, stakeholder, iter))
	if errQr != nil {
		return [](*taps.Threadrel){}, fmt.Errorf("Could not query for threads: %v", errQr)
	}
	defer qr.Close()
	ths := [](*taps.Threadrel){}
	for qr.Next() {
		th := &taps.Threadrel{}
		errScn := qr.Scan(&th.ID, &th.State, &th.Owner, &th.Percentile, &th.Iteration, &th.CostDirect, &th.Order)
		if errScn != nil {
			return [](*taps.Threadrel){}, fmt.Errorf("Could not scan thread: %v", errScn)
		}
		th.StakeholderMatch = true
		ths = append(ths, th)
	}
	return ths, nil
}
*/

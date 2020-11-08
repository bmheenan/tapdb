package tapdb

/*

import (
	"database/sql"
	"fmt"

	"github.com/bmheenan/taps"
)

// GetThreadrel returns a Threadrel for the matching thread `id`. `StakeholderMatch` and `Order` will not be filled
func (db *mysqlDB) GetThreadrel(id int64, stakeholder string) (*taps.Threadrel, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
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
	FROM        threads t
	  LEFT JOIN (
				SELECT thread
				  ,    stakeholder
				  ,    ord
				FROM   threads_stakeholders
				WHERE  stakeholder = '%v'
				  AND  thread = %v
				) AS s
	  ON        t.id = s.thread
	WHERE       id = %v
	;`, stakeholder, stakeholder, id, id))
	if errQry != nil {
		return &taps.Threadrel{}, fmt.Errorf("Could not query for thread: %v", errQry)
	}
	defer qr.Close()
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
			return &taps.Threadrel{}, fmt.Errorf("Could not scan thread: %v", errScn)
		}
		return th, nil
	}
	return &taps.Threadrel{}, fmt.Errorf("No thread found with id %v: %w", id, ErrNotFound)
}

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

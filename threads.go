package tapdb

import (
	"fmt"

	"github.com/bmheenan/taps"
)

// NewThread inserts a new thread into the db with the given data. It assumes the thread has no children.
// Returns the id of the newly inserted thread or an error
func (db *mysqlDB) NewThread(name, domain, owner, iteration, state string, percentile float64, cost int) (int64, error) {
	if name == "" || domain == "" || owner == "" || iteration == "" || state == "" || percentile < 0 || cost < 0 {
		return 0, fmt.Errorf("Args must be non-blank; cost and percenitle must be >= 0: %w", ErrBadArgs)
	}
	res, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads
	            (name, domain, owner, iteration, state, percentile, costdirect, costtotal)
	VALUES      ('%v',   '%v',  '%v',      '%v',  '%v',         %v,         %v,        %v)
	;`, name, domain, owner, iteration, state, percentile, cost, cost))
	if errIn != nil {
		return 0, fmt.Errorf("Could not insert new thread into db: %v", errIn)
	}
	id, errID := res.LastInsertId()
	if errID != nil {
		return 0, fmt.Errorf("Could not get id of inserted thread: %v", errID)
	}
	return id, nil
}

func (db *mysqlDB) LinkThreads(parent, child int64, iter string, ord int, domain string) error {
	if ord < 0 || domain == "" || iter == "" {
		return fmt.Errorf("Domain and iteration must be non-blank; order must be >= 0: %w", ErrBadArgs)
	}
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_parent_child
	            (parent, child, domain, iteration, ord)
	VALUES      (    %v,    %v,   '%v',      '%v',  %v)
	;`, parent, child, domain, iter, ord))
	return err
}

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

func (db *mysqlDB) GetThreadOrderBefore(parent int64, iter string, order int) (int, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_parent_child
	WHERE  parent = %v
	  AND  ord < %v
	  AND  iteration = '%v'
	;`, parent, order, iter))
	if errQry != nil {
		return 0, fmt.Errorf("Could not query for previous thread order: %v", errQry)
	}
	defer qr.Close()
	max := 0
	for qr.Next() {
		errScn := qr.Scan(&max)
		if errScn != nil {
			return 0, fmt.Errorf("Could not scan max value: %v", errScn)
		}
	}
	return max + ((order - max) / 2), nil
}

//func (db *mysqlDB) GetLowestAncestors()

/*
	WITH        RECURSIVE descendants (child, parent) AS
	            (
	            SELECT child
	              ,    parent
	            FROM   threads_parent_child
	            WHERE  parent = '%v'
	            UNION ALL
	            SELECT t.child
	              ,    t.parent
	            FROM   threads_parent_child t
	            JOIN   descendants
	              ON   t.parent = descendants.child
		        )
	SELECT      d.child
	  ,         t.owner
	  ,         (s.stakeholder = '%v') AS tracked
	  , 	    t.costdirect
	  , 	    t.iteration
	FROM        descendants d
	  LEFT JOIN (
			    SELECT thread
			      ,    stakeholder
				FROM   threads_stakeholders
				WHERE  stakeholder = '%v'
	            ) s
	  ON        s.thread = d.child
	  JOIN      threads t
	  ON        t.id = d.child
	ORDER BY    s.stakeholder
*/

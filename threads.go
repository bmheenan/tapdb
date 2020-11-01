package tapdb

import (
	"fmt"
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

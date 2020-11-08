package tapdb

/*
import (
	"errors"
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

func (db *mysqlDB) LinkThreadsStakeholder(parent, child int64, stakeholder, domain string) error {
	if domain == "" || stakeholder == "" {
		return errors.New("Doman and stakeholder must be non-blank")
	}
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders_parent_child
				(parent, child, stakeholder, domain)
	VALUES      (    %v,    %v,        '%v',   '%v')
	;`, parent, child, stakeholder, domain))
	return err
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
			return 0, nil
		}
	}
	return max + ((order - max) / 2), nil
}

func (db *mysqlDB) GetPersonteamOrderBefore(personteam, iter string, order int) (int, error) {
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT MAX(ord) AS ord
	FROM   threads_stakeholders
	WHERE  stakeholder = %v
	  AND  ord < %v
	  AND  iteration = '%v'
	;`, personteam, order, iter))
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
	return max + ((order - max) / 2), nil
}

func (db *mysqlDB) SetThreadCostTotal(id int64, cost int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads
	SET    costtotal = %v
	WHERE  id = %v
	;`, cost, id))
	return err
}

func (db *mysqlDB) SetThreadOrderParent(thread, parent int64, order int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_parent_child
	SET    ord = %v
	WHERE  child = %v
	  AND  parent = %v
	;`, order, thread, parent))
	return err
}

func (db *mysqlDB) SetThreadOrderStakeholder(thread int64, stakeholder string, order int) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	UPDATE threads_stakeholders
	SET    ord = %v
	WHERE  thread = %v
	  AND  stakeholder = %v
	;`, order, thread, stakeholder))
	return err
}
*/

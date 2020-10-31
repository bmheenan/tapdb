package tapdb

import (
	"errors"
	"fmt"
)

const keyNewThreadParentLink = "newthreadparentlink"
const qryNewThreadParentLink = `
INSERT INTO threads_parent_child (parent, child, domain, ord)
VALUES                           (     ?,     ?,      ?,   ?);`

func (db *mySQLDB) initThreadLink() error {
	var err error
	db.stmts[keyNewThreadParentLink], err = db.conn.Prepare(qryNewThreadParentLink)
	if err != nil {
		return fmt.Errorf("Could not init %v: %v", keyNewThreadParentLink, err)
	}
	return nil
}

func (db *mySQLDB) ThreadLink(parent int64, child int64) error {
	return errors.New("Not implemented")
}

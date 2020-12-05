package tapdb

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	_, err := Init(getTestCredentials())
	if err != nil {
		t.Errorf("Init returned error: %v", err)
		return
	}
}

func setupEmptyDB() (DBInterface, error) {
	db, errS := Init(getTestCredentials())
	if errS != nil {
		return &mysqlDB{}, errS
	}
	db.ClearThreadStkLinks("example.com")
	db.ClearThreadHierLinks("example.com")
	db.ClearThreads("example.com")
	db.ClearStkHierLinks("example.com")
	db.ClearStks("example.com")
	return db, nil
}

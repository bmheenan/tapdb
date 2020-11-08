package tapdb

import (
	"fmt"
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
	errThStkH := db.ClearThreadStkHierLinks("example.com")
	if errThStkH != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear threads/stakeholders heirarchy links: %v", errThStkH)
	}
	errThStk := db.ClearThreadStkLinks("example.com")
	if errThStk != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear threads/stakeholders links: %v", errThStk)
	}
	errThH := db.ClearThreadHierLinks("example.com")
	if errThH != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear threads heirarchy links: %v", errThH)
	}
	errTh := db.ClearThreads("example.com")
	if errTh != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear threads: %v", errTh)
	}
	errStkH := db.ClearStkHierLinks("example.com")
	if errStkH != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear stakeholders heirarchy: %v", errStkH)
	}
	errStk := db.ClearStks("example.com")
	if errStk != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear stakeholders: %v", errStk)
	}
	return db, nil
}

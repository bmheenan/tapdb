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

func setupForTest() (DBInterface, error) {
	db, err := Init(getTestCredentials())
	if err != nil {
		return &mysqlDB{}, fmt.Errorf("Init returned error: %v", err)
	}
	errSk := db.ClearStakeholders("example.com")
	if errSk != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear stakeholders: %v", errSk)
	}
	errTPC := db.ClearThreadsPC("example.com")
	if errTPC != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear thread parent child relationships: %v", errTPC)
	}
	errTh := db.ClearThreads("example.com")
	if errTh != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear threads: %v", errTh)
	}
	errCPC := db.ClearPersonteamsPC("example.com")
	if errCPC != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear personteam parent/child relationships: %v", errCPC)
	}
	errCPT := db.ClearPersonteams("example.com")
	if errCPT != nil {
		return &mysqlDB{}, fmt.Errorf("Could not clear personteams: %v", errCPT)
	}
	return db, nil
}

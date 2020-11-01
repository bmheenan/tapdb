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

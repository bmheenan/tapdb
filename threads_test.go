package tapdb

import (
	"fmt"
	"testing"
)

func TestNewAndGetThread(t *testing.T) {
	db, pts, errSetup := setupWithPersonteams()
	if errSetup != nil {
		t.Errorf("%v", errSetup)
		return
	}
	id, errIn := db.NewThread("Test thread", "example.com", pts[0], "2020 Oct", "not started", 1, 10)
	if errIn != nil {
		t.Errorf("Could not insert new thread: %v", errIn)
		return
	}
	th, errGet := db.GetThreadrel(id)
	if errGet != nil {
		t.Errorf("Could not get thread: %v", errGet)
		return
	}
	if th.ID != id || th.Iteration != "2020 Oct" || th.CostDirect != 10 {
		t.Errorf("Retrieved thread didn't have the expected data. Got %v", th)
		return
	}
}

func setupWithPersonteams() (DBInterface, []string, error) {
	db, errSetup := setupForTest()
	if errSetup != nil {
		return nil, nil, fmt.Errorf("Could not set up test: %v", errSetup)
	}
	es := []string{
		"a@example.com",
		"b@example.com",
		"c@example.com",
	}
	for _, e := range es {
		errPT := db.NewPersonteam(e, "example.com", "PT", "PT", "#ffffff", "#000000", "monthly")
		if errPT != nil {
			return nil, nil, fmt.Errorf("Error trying to insert new personteam: %v", errPT)
		}
	}
	return db, es, nil
}

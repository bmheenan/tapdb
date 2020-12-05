package tapdb

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestGetItersForStk(t *testing.T) {
	db, stks := setupWithStks()
	id1 := db.NewThread("A", "example.com", stks[0], "2020 Q4", string(taps.NotStarted), 1, 1)
	db.NewThreadStkLink(id1, stks[0], "example.com", "2020 Q4", 1, 1)
	id2 := db.NewThread("B", "example.com", stks[0], "2021 Q1", string(taps.NotStarted), 1, 1)
	db.NewThreadStkLink(id2, stks[0], "example.com", "2021 Q1", 1, 1)
	iters, errI := db.GetItersForStk(stks[0])
	if errI != nil {
		t.Errorf("Could not get iterations: %v", errI)
		return
	}
	if len(iters) != 2 {
		t.Errorf("Expected 2 iterations; got %v", len(iters))
		return
	}
	if iters[0] != "2020 Q4" {
		t.Errorf("Expected first result to be 2020 Q4; got %v", iters[0])
		return
	}
	if iters[1] != "2021 Q1" {
		t.Errorf("Expected second result to be 2021 Q1; got %v", iters[1])
		return
	}
}

func TestGetItersForParent(t *testing.T) {
	db, stks := setupWithStks()
	idA := db.NewThread("A", "example.com", stks[0], "2020 Dec", string(taps.NotStarted), 1, 1)
	idAA := db.NewThread("AA", "example.com", stks[0], "2020-10 Oct", string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(idA, idAA, "2020-10 Oct", 1, "example.com")
	idAB := db.NewThread("AB", "example.com", stks[0], "2020-11 Nov", string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(idA, idAB, "2020-11 Nov", 1, "example.com")
	idAC := db.NewThread("AC", "example.com", stks[0], "2020-12 Dec", string(taps.NotStarted), 1, 1)
	db.NewThreadHierLink(idA, idAC, "2020-12 Dec", 1, "example.com")
	iters, errI := db.GetItersForParent(idA)
	if errI != nil {
		t.Errorf("Could not get iterations: %v", errI)
		return
	}
	if len(iters) != 3 {
		t.Errorf("Expected 3 iterations; got %v", len(iters))
		return
	}
	if iters[0] != "2020-10 Oct" {
		t.Errorf("Expected first result to be 2020 Oct; got %v", iters[0])
		return
	}
	if iters[1] != "2020-11 Nov" {
		t.Errorf("Expected second result to be 2021 Nov; got %v", iters[1])
		return
	}
	if iters[2] != "2020-12 Dec" {
		t.Errorf("Expected second result to be 2021 Dec; got %v", iters[1])
		return
	}
}

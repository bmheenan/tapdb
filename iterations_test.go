package tapdb

import (
	"testing"

	"github.com/bmheenan/taps"
)

func TestGetItersForStk(t *testing.T) {
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	id1, errN1 := db.NewThread("A", "example.com", stks[0], "2020 Q4", string(taps.NotStarted), 1, 1)
	if errN1 != nil {
		t.Errorf("Could not insert thread A: %v", errN1)
		return
	}
	errL1 := db.NewThreadStkLink(id1, stks[0], "example.com", "2020 Q4", 1, true, 1)
	if errL1 != nil {
		t.Errorf("Could not link thread to stakeholder: %v", errL1)
		return
	}
	id2, errN2 := db.NewThread("B", "example.com", stks[0], "2021 Q1", string(taps.NotStarted), 1, 1)
	if errN2 != nil {
		t.Errorf("Could not insert thread B: %v", errN2)
		return
	}
	errL2 := db.NewThreadStkLink(id2, stks[0], "example.com", "2021 Q1", 1, true, 1)
	if errL2 != nil {
		t.Errorf("Could not link thread to stakeholder: %v", errL2)
		return
	}
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
	db, stks, errSetup := setupWithStks()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	idA, errNA := db.NewThread("A", "example.com", stks[0], "2020 Dec", string(taps.NotStarted), 1, 1)
	if errNA != nil {
		t.Errorf("Could not insert thread A: %v", errNA)
		return
	}
	idAA, errNAA := db.NewThread("AA", "example.com", stks[0], "2020-10 Oct", string(taps.NotStarted), 1, 1)
	if errNAA != nil {
		t.Errorf("Could not insert thread AA: %v", errNAA)
		return
	}
	errLAA := db.NewThreadHierLink(idA, idAA, "2020-10 Oct", 1, "example.com")
	if errLAA != nil {
		t.Errorf("Could not link thread AA with A: %v", errLAA)
		return
	}
	idAB, errNAB := db.NewThread("AB", "example.com", stks[0], "2020-11 Nov", string(taps.NotStarted), 1, 1)
	if errNAB != nil {
		t.Errorf("Could not insert thread AB: %v", errNAB)
		return
	}
	errLAB := db.NewThreadHierLink(idA, idAB, "2020-11 Nov", 1, "example.com")
	if errLAB != nil {
		t.Errorf("Could not link thread AB with A: %v", errLAB)
		return
	}
	idAC, errNAC := db.NewThread("AC", "example.com", stks[0], "2020-12 Dec", string(taps.NotStarted), 1, 1)
	if errNAC != nil {
		t.Errorf("Could not insert thread AC: %v", errNAC)
		return
	}
	errLAC := db.NewThreadHierLink(idA, idAC, "2020-12 Dec", 1, "example.com")
	if errLAC != nil {
		t.Errorf("Could not link thread AC with A: %v", errLAC)
		return
	}
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

package tapdb

import (
	"errors"
	"testing"
)

func TestNewAndGetStk(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	errNew := db.NewStk("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	if errNew != nil {
		t.Errorf("Error trying to insert new stakeholder: %v", errNew)
	}
	stk, errGet := db.GetStk("a@example.com")
	if errGet != nil {
		t.Errorf("GetPersonteam returned an error: %v", errGet)
		return
	}
	if stk.Email != "a@example.com" {
		t.Errorf("Expected email to be a@example.com but got %v", stk.Email)
	}
}

func TestGetStkNotFound(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	errNew := db.NewStk("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	if errNew != nil {
		t.Errorf("Error trying to insert new stakeholder: %v", errNew)
	}
	pt, errGet := db.GetStk("b@example.com")
	if errGet == nil {
		t.Errorf("GetStk didn't return an error for a stakeholder that didn't exist. Returned: %v", pt)
		return
	}
	if !errors.Is(errGet, ErrNotFound) {
		t.Errorf("Returned error was not ErrNotFound")
		return
	}
}

func TestGetStkDes(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	allStks := []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
		"b@example.com",
		"ba@example.com",
	}
	for _, e := range allStks {
		errNew := db.NewStk(e, "example.com", "Personteam", "A", "#ffffff", "#000000", "monthly")
		if errNew != nil {
			t.Errorf("Error trying to insert new stakeholder %v: %v", e, errNew)
			return
		}
	}
	for _, l := range []struct {
		p string
		c string
	}{
		{
			p: "a@example.com",
			c: "aa@example.com",
		},
		{
			p: "aa@example.com",
			c: "aaa@example.com",
		},
		{
			p: "a@example.com",
			c: "ab@example.com",
		},
		{
			p: "b@example.com",
			c: "ba@example.com",
		},
	} {
		errPC := db.NewStkHierLink(l.p, l.c, "example.com")
		if errPC != nil {
			t.Errorf("Error trying to make %v a parent of %v: %v", l.p, l.c, errPC)
			return
		}
	}
	stks0, errG0 := db.GetStkDes("a@example.com")
	if errG0 != nil {
		t.Errorf("GetStkDes returned an error: %v", errG0)
		return
	}
	if len(stks0) != 4 {
		t.Errorf("GetStkDes expected length 4, but got %v", len(stks0))
		return
	}
	for _, e := range []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
	} {
		if _, ok := stks0[e]; !ok {
			t.Errorf("GetStkDes (0) was missing %v", e)
			return
		}
	}
	stks1, errG1 := db.GetStkDes("b@example.com")
	if errG1 != nil {
		t.Errorf("GetStkDes returned an error: %v", errG1)
		return
	}
	if len(stks1) != 2 {
		t.Errorf("GetStkDes expected length 2, got %v", len(stks1))
	}
	for _, e := range []string{"b@example.com", "ba@example.com"} {
		if _, ok := stks1[e]; !ok {
			t.Errorf("GetStkDes (1) was missing %v", e)
			return
		}
	}
}

func TestGetStkAns(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	allStks := []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
		"b@example.com",
		"ba@example.com",
	}
	for _, e := range allStks {
		errNew := db.NewStk(e, "example.com", "Personteam", "A", "#ffffff", "#000000", "monthly")
		if errNew != nil {
			t.Errorf("Error trying to insert new stakeholder %v: %v", e, errNew)
			return
		}
	}
	for _, l := range []struct {
		p string
		c string
	}{
		{
			p: "a@example.com",
			c: "aa@example.com",
		},
		{
			p: "aa@example.com",
			c: "aaa@example.com",
		},
		{
			p: "a@example.com",
			c: "ab@example.com",
		},
		{
			p: "b@example.com",
			c: "ba@example.com",
		},
	} {
		errPC := db.NewStkHierLink(l.p, l.c, "example.com")
		if errPC != nil {
			t.Errorf("Error trying to make %v a parent of %v: %v", l.p, l.c, errPC)
			return
		}
	}
	stks0, errG0 := db.GetStkAns("aaa@example.com")
	if errG0 != nil {
		t.Errorf("GetStkAns returned an error: %v", errG0)
		return
	}
	if len(stks0) != 3 {
		t.Errorf("GetStkAns expected length 3, but got %v", len(stks0))
		return
	}
	for _, e := range []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
	} {
		if _, ok := stks0[e]; !ok {
			t.Errorf("GetStkAns (0) was missing %v", e)
			return
		}
	}
	stks1, errG1 := db.GetStkAns("b@example.com")
	if errG1 != nil {
		t.Errorf("GetStkAns returned an error: %v", errG1)
		return
	}
	if len(stks1) != 1 {
		t.Errorf("GetStkAns expected length 1, got %v", len(stks1))
	}
	for _, e := range []string{"b@example.com"} {
		if _, ok := stks1[e]; !ok {
			t.Errorf("GetStkDes (1) was missing %v", e)
			return
		}
	}
}

func TestPTDescendantsNotFound(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	errNew := db.NewStk("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	if errNew != nil {
		t.Errorf("Error trying to insert new personteam: %v", errNew)
	}
	stks, errGet := db.GetStkDes("b@example.com")
	if errGet == nil {
		t.Errorf("GetStkDes didn't return an error for a stakeholder that didn't exist. Returned: %v", stks)
		return
	}
	if !errors.Is(errGet, ErrNotFound) {
		t.Errorf("Returned error was not ErrNotFound")
		return
	}
}

func TestGetTeamsForDomain(t *testing.T) {
	db, errSetup := setupEmptyDB()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	allStks := []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
		"b@example.com",
		"ba@example.com",
	}
	for _, e := range allStks {
		errNew := db.NewStk(e, "example.com", "Personteam", "A", "#ffffff", "#000000", "monthly")
		if errNew != nil {
			t.Errorf("Error trying to insert new stakeholder %v: %v", e, errNew)
			return
		}
	}
	for _, l := range []struct {
		p string
		c string
	}{
		{
			p: "a@example.com",
			c: "aa@example.com",
		},
		{
			p: "aa@example.com",
			c: "aaa@example.com",
		},
		{
			p: "a@example.com",
			c: "ab@example.com",
		},
		{
			p: "b@example.com",
			c: "ba@example.com",
		},
	} {
		errPC := db.NewStkHierLink(l.p, l.c, "example.com")
		if errPC != nil {
			t.Errorf("Error trying to make %v a parent of %v: %v", l.p, l.c, errPC)
			return
		}
	}
	ret, errR := db.GetStksForDomain("example.com")
	if errR != nil {
		t.Errorf("Could not get all stakeholders: %v", errR)
		return
	}
	if len(ret) != 2 {
		t.Errorf("Expected 2 results; got %v", len(ret))
		return
	}
	if ret[0].Email != "a@example.com" {
		t.Errorf("expected a; got %v", ret[0].Email)
		return
	}
	if ret[0].Members[0].Email != "aa@example.com" {
		t.Errorf("expected aa; got %v", ret[0].Email)
		return
	}
	if ret[0].Members[0].Members[0].Email != "aaa@example.com" {
		t.Errorf("expected aaa; got %v", ret[0].Email)
		return
	}
	if ret[0].Members[1].Email != "ab@example.com" {
		t.Errorf("expected ab; got %v", ret[0].Email)
		return
	}
	if ret[1].Email != "b@example.com" {
		t.Errorf("expected b; got %v", ret[0].Email)
		return
	}
	if ret[1].Members[0].Email != "ba@example.com" {
		t.Errorf("expected ba; got %v", ret[0].Email)
		return
	}
}

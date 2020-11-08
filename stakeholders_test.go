package tapdb

import (
	"testing"
)

func TestNewAndGetPersonteam(t *testing.T) {
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

/*
func TestPTNotFound(t *testing.T) {
	db, errSetup := setupForTest()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	errNew := db.NewPersonteam("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	if errNew != nil {
		t.Errorf("Error trying to insert new personteam: %v", errNew)
	}
	pt, errGet := db.GetPersonteam("b@example.com")
	if errGet == nil {
		t.Errorf("GetPersonteam didn't return an error for a personteam that didn't exist. Returned: %v", pt)
		return
	}
	if !errors.Is(errGet, ErrNotFound) {
		t.Errorf("Returned error was not ErrNotFound")
		return
	}
}

func TestGetPersonteamDescendants(t *testing.T) {
	db, errSetup := setupForTest()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	allPTs := []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
		"b@example.com",
		"ba@example.com",
	}
	for _, e := range allPTs {
		errNew := db.NewPersonteam(e, "example.com", "Personteam", "A", "#ffffff", "#000000", "monthly")
		if errNew != nil {
			t.Errorf("Error trying to insert new personteam %v: %v", e, errNew)
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
		errPC := db.LinkPersonteams(l.p, l.c, "example.com")
		if errPC != nil {
			t.Errorf("Error trying to link personteams %v and %v: %v", l.p, l.c, errPC)
			return
		}
	}
	pts0, errG0 := db.GetPersonteamDescendants("a@example.com")
	if errG0 != nil {
		t.Errorf("GetPersonteamDescendants returned an error: %v", errG0)
		return
	}
	if len(pts0) != 4 {
		t.Errorf("GetpersonteamDescendants expected length %v, got %v", len(allPTs), 4)
	}
	for _, e := range []string{
		"a@example.com",
		"aa@example.com",
		"aaa@example.com",
		"ab@example.com",
	} {
		if _, ok := pts0[e]; !ok {
			t.Errorf("GetPersonteamDescendants was missing %v", e)
			return
		}
	}
	pts1, errG1 := db.GetPersonteamDescendants("b@example.com")
	if errG1 != nil {
		t.Errorf("GetPersonteamDescendants returned an error: %v", errG1)
		return
	}
	if len(pts1) != 2 {
		t.Errorf("GetpersonteamDescendants expected length 2, got %v", len(pts1))
	}
	for _, e := range []string{"b@example.com", "ba@example.com"} {
		if _, ok := pts1[e]; !ok {
			t.Errorf("GetPersonteamDescendants was missing %v", e)
			return
		}
	}
}

func TestPTDescendantsNotFound(t *testing.T) {
	db, errSetup := setupForTest()
	if errSetup != nil {
		t.Errorf("Could not set up test: %v", errSetup)
		return
	}
	errNew := db.NewPersonteam("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	if errNew != nil {
		t.Errorf("Error trying to insert new personteam: %v", errNew)
	}
	pt, errGet := db.GetPersonteamDescendants("b@example.com")
	if errGet == nil {
		t.Errorf("GetPersonteamDescendants didn't return an error for a personteam that didn't exist. Returned: %v", pt)
		return
	}
	if !errors.Is(errGet, ErrNotFound) {
		t.Errorf("Returned error was not ErrNotFound")
		return
	}
}
*/

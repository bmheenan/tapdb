package tapdb

import "testing"

func TestNewAndGetPersonteam(t *testing.T) {
	db, errInit := Init(getTestCredentials())
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	db.ClearPersonteams("example.com")
	db.NewPersonteam("a@example.com", "example.com", "Team A", "A", "#ffffff", "#000000", "monthly")
	pt, errGet := db.GetPersonteam("a@example.com")
	if errGet != nil {
		t.Errorf("GetPersonteam returned an error: %v", errGet)
		return
	}
	if pt.Email != "a@example.com" {
		t.Errorf("Expected email to be a@example.com but got %v", pt.Email)
	}
}

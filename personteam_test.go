package tapdb

import (
	"errors"
	"testing"

	"github.com/bmheenan/tapstruct"
)

func TestNewAndGetPersonteam(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Clear domain returned error: %v", errClear)
		return
	}
	errNew := db.NewPersonteam(&tapstruct.Personteam{
		Email:      "adam@example.com",
		Domain:     "example.com",
		Name:       "Adam",
		Abbrev:     "AD",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}, "")
	if errNew != nil {
		t.Errorf("NewPersonteam returned an error: %v", errNew)
		return
	}
	pt, errGet := db.GetPersonteam("adam@example.com", 0)
	if errGet != nil {
		t.Errorf("GetPersonteam returned an error: %v", errGet)
		return
	}
	if pt.Email != "adam@example.com" ||
		pt.Domain != "example.com" ||
		pt.Name != "Adam" ||
		pt.Abbrev != "AD" ||
		pt.ColorF != "#ffffff" ||
		pt.ColorB != "#1f57cf" ||
		pt.IterTiming != tapstruct.Monthly {
		t.Errorf("GetPersonteam didn't have the expected results: %v", pt)
		return
	}
}

func TestNewBulkAndGetPersonteam(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Clear domain returned error: %v", errClear)
		return
	}
	errNew := db.NewPersonteam(&tapstruct.Personteam{
		Email:       "team@example.com",
		Domain:      "example.com",
		Name:        "Team",
		Abbrev:      "TM",
		ColorF:      "#ffffff",
		ColorB:      "#1f57cf",
		IterTiming:  tapstruct.Monthly,
		HasChildren: true,
		Children: []tapstruct.Personteam{
			tapstruct.Personteam{
				Email:      "taylor@example.com",
				Domain:     "example.com",
				Name:       "Taylor",
				Abbrev:     "TA",
				ColorF:     "#ffffff",
				ColorB:     "#1f57cf",
				IterTiming: tapstruct.Biweekly,
			},
		},
	}, "")
	if errNew != nil {
		t.Errorf("NewPersonteam returned an error: %v", errNew)
		return
	}
	pt, errGet := db.GetPersonteam("team@example.com", 1)
	if errGet != nil {
		t.Errorf("GetPersonteam returned an error: %v", errGet)
		return
	}
	if pt.Email != "team@example.com" ||
		pt.Domain != "example.com" ||
		pt.Name != "Team" ||
		pt.Abbrev != "TM" ||
		pt.ColorF != "#ffffff" ||
		pt.ColorB != "#1f57cf" ||
		pt.IterTiming != tapstruct.Monthly {
		t.Errorf("GetPersonteam didn't have the expected results: %v", pt)
		return
	}
	if !pt.HasChildren || len(pt.Children) != 1 {
		t.Errorf("The team didn't have one child")
		return
	}
	child := pt.Children[0]
	if child.Email != "taylor@example.com" ||
		child.Name != "Taylor" ||
		child.HasChildren != false ||
		child.IterTiming != tapstruct.Biweekly {
		t.Errorf("The child didn't have the expected results: %v", child)
		return
	}
}

func TestNewAndGetPersonteamTree(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Clear domain returned error: %v", errClear)
		return
	}
	errNew := db.NewPersonteam(&tapstruct.Personteam{
		Email:      "team@example.com",
		Domain:     "example.com",
		Name:       "Team",
		Abbrev:     "TM",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}, "")
	if errNew != nil {
		t.Errorf("NewPersonteam returned an error: %v", errNew)
		return
	}
	errNew = db.NewPersonteam(&tapstruct.Personteam{
		Email:      "eve@example.com",
		Domain:     "example.com",
		Name:       "Eve",
		Abbrev:     "EVE",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Biweekly,
	}, "team@example.com")
	if errNew != nil {
		t.Errorf("NewPersonteam returned an error when inserting child: %v", errNew)
		return
	}
	pt, errGet := db.GetPersonteam("team@example.com", 1)
	if errGet != nil {
		t.Errorf("GetPersonteam returned an error: %v", errGet)
		return
	}
	if pt.Email != "team@example.com" ||
		pt.Domain != "example.com" ||
		pt.Name != "Team" ||
		pt.Abbrev != "TM" ||
		pt.ColorF != "#ffffff" ||
		pt.ColorB != "#1f57cf" ||
		pt.IterTiming != tapstruct.Monthly {
		t.Errorf("GetPersonteam didn't have the expected results: %v", pt)
		return
	}
	if !pt.HasChildren || len(pt.Children) != 1 {
		t.Errorf("The team didn't have one child")
		return
	}
	child := pt.Children[0]
	if child.Email != "eve@example.com" ||
		child.Name != "Eve" ||
		child.HasChildren != false ||
		child.IterTiming != tapstruct.Biweekly {
		t.Errorf("The child didn't have the expected results: %v", child)
		return
	}
}

func TestGetNonExistingPersonteam(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	errClear := db.ClearDomain("example.com")
	if errClear != nil {
		t.Errorf("Clear domain returned error: %v", errClear)
		return
	}
	errNew := db.NewPersonteam(&tapstruct.Personteam{
		Email:      "team@example.com",
		Domain:     "example.com",
		Name:       "Team",
		Abbrev:     "TM",
		ColorF:     "#ffffff",
		ColorB:     "#1f57cf",
		IterTiming: tapstruct.Monthly,
	}, "")
	if errNew != nil {
		t.Errorf("NewPersonteam returned an error: %v", errNew)
		return
	}
	pt, errGet := db.GetPersonteam("eve@example.com", 0)
	if errGet == nil {
		t.Errorf("Found a result that shouldn't have existed: %v", pt)
		return
	} else if !errors.Is(errGet, ErrNotFound) {
		t.Errorf("Found an error getting a result that doesn't exist, but it's not the right error: %v", errGet)
	}
}

func TestGetPersonteamDepthTooHigh(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	_, err := db.GetPersonteam("ex@example.com", 6)
	if err == nil {
		t.Errorf("Did not receive an error when GetPersonteam's depth was too high")
	}
}

func TestGetPersonteamBlankEmail(t *testing.T) {
	db, errInit := InitDB()
	if errInit != nil {
		t.Errorf("Init returned error: %v", errInit)
		return
	}
	_, err := db.GetPersonteam("", 0)
	if err == nil {
		t.Errorf("Did not receive an error when GetPersonteam's email was blank")
	}
}

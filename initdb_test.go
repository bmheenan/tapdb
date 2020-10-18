package tapdb

import (
	"testing"
)

func TestInitDb(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Errorf("Init returned error: %v", err)
		return
	}
	if m, ok := db.(*mySQLDB); ok {
		_, errDescPT := m.conn.Exec("DESCRIBE personteams")
		if errDescPT != nil {
			t.Errorf("Could not verify personteams exists: %v", errDescPT)
		}
	}
}

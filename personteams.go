package tapdb

import (
	"fmt"

	taps "github.com/bmheenan/taps"
)

// NewPersonteam inserts a new personteam into the db with the given data
func (db *mysqlDB) NewPersonteam(email, domain, name, abbrev, colorf, colorb string, itertiming taps.IterTiming) error {
	if email == "" || domain == "" || name == "" || abbrev == "" || colorb == "" || colorf == "" {
		return fmt.Errorf("No args may be blank: %w", ErrBadArgs)
	}
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO personteams
	            (email, domain, name, abbrev, colorf, colorb, itertiming)
	VALUES      ( '%v',   '%v', '%v',   '%v',   '%v',   '%v',       '%v')
	;`, email, domain, name, abbrev, colorf, colorb, itertiming))
	return err
}

// NewPersonteamPC makes `child` a child of `parent`. Both must already exist
func (db *mysqlDB) LinkPersonteams(parent, child, domain string) error {
	if parent == "" || child == "" || domain == "" {
		return fmt.Errorf("No args may be blank: %w", ErrBadArgs)
	}
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO personteams_parent_child
	            (parent, child, domain)
	VALUES      (  '%v',  '%v',   '%v')
	;`, parent, child, domain))
	return err
}

// GetPersonteam gets the details for the given personteam, without details of any children
func (db *mysqlDB) GetPersonteam(email string) (*taps.Personteam, error) {
	if email == "" {
		return &taps.Personteam{}, fmt.Errorf("Email cannot be blank: %w", ErrBadArgs)
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	SELECT  email
	  ,     domain
	  ,     name
	  ,     abbrev
	  ,     colorf
	  ,     colorb
	  ,     itertiming
	FROM    personteams
	WHERE   email = '%v'
	;`, email))
	if errQry != nil {
		return &taps.Personteam{}, fmt.Errorf("Could not query for personteam %v: %v", email, errQry)
	}
	defer qr.Close()
	pt := &taps.Personteam{}
	found := false
	for qr.Next() {
		found = true
		errScn := qr.Scan(&pt.Email, &pt.Domain, &pt.Name, &pt.Abbrev, &pt.ColorF, &pt.ColorB, &pt.IterTiming)
		if errScn != nil {
			return &taps.Personteam{}, fmt.Errorf("Could not scan personteam: %v", errScn)
		}
	}
	if !found {
		return pt, fmt.Errorf("No personteam by that email: %w", ErrNotFound)
	}
	return pt, nil
}

// GetPersonteamDescendants gets a map of all personteams that are descendants of the given one (including itself).
// The hierarchy is not preserved.
func (db *mysqlDB) GetPersonteamDescendants(email string) (map[string](*taps.Personteam), error) {
	if email == "" {
		return map[string](*taps.Personteam){}, fmt.Errorf("Email cannot be blank: %w", ErrBadArgs)
	}
	topPT, errPT := db.GetPersonteam(email)
	if errPT != nil {
		return map[string](*taps.Personteam){}, fmt.Errorf("Could not get top level personteam %v: %w", email, errPT)
	}
	pts := map[string](*taps.Personteam){
		email: topPT,
	}
	qr, errQry := db.conn.Query(fmt.Sprintf(`
	WITH   RECURSIVE descendants (child, parent) AS
	       (
	       SELECT child
	         ,    parent
	       FROM   personteams_parent_child
	       WHERE  parent = '%v'
	       UNION ALL
	       SELECT pt.child
	         ,    pt.parent
	       FROM   personteams_parent_child pt
	       JOIN   descendants d
	         ON   pt.parent = d.child
	       )
	SELECT pt.email
	  ,    pt.domain
	  ,    pt.name
	  ,    pt.abbrev
	  ,    pt.colorf
	  ,    pt.colorb
	  ,    pt.itertiming
	FROM   descendants d
	  JOIN personteams pt
	  ON   pt.email = d.child
	;`, email))
	if errQry != nil {
		return map[string](*taps.Personteam){}, fmt.Errorf("Could not query for personteam %v: %v", email, errQry)
	}
	defer qr.Close()
	for qr.Next() {
		pt := taps.Personteam{}
		errScn := qr.Scan(&pt.Email, &pt.Domain, &pt.Name, &pt.Abbrev, &pt.ColorF, &pt.ColorB, &pt.IterTiming)
		if errScn != nil {
			return map[string](*taps.Personteam){}, fmt.Errorf("Could not scan personteam: %v", errScn)
		}
		pts[pt.Email] = &pt
	}
	return pts, nil
}

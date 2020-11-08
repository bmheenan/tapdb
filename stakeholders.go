package tapdb

import (
	"fmt"

	taps "github.com/bmheenan/taps"
)

func (db *mysqlDB) NewStk(email, domain, name, abbrev, colorf, colorb string, cadence taps.Cadence) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO stakeholders
	            (email, domain, name, abbrev, colorf, colorb, cadence)
	VALUES      ( '%v',   '%v', '%v',   '%v',   '%v',   '%v',    '%v')
	;`, email, domain, name, abbrev, colorf, colorb, string(cadence)))
	return err
}

func (db *mysqlDB) NewStkHierLink(parent, child, domain string) error {
	_, err := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO stakeholders_hierarchy
	            (parent, child, domain)
	VALUES      (  '%v',  '%v',   '%v')
	;`, parent, child, domain))
	return err
}

func (db *mysqlDB) GetStk(email string) (*taps.Stakeholder, error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT  email
	  ,     domain
	  ,     name
	  ,     abbrev
	  ,     colorf
	  ,     colorb
	  ,     cadence
	FROM    stakeholders
	WHERE   email = '%v'
	;`, email))
	if errQr != nil {
		return &taps.Stakeholder{}, fmt.Errorf("Could not query for stakeholder %v: %v", email, errQr)
	}
	defer qr.Close()
	stk := &taps.Stakeholder{}
	found := false
	for qr.Next() {
		found = true
		errScn := qr.Scan(&stk.Email, &stk.Domain, &stk.Name, &stk.Abbrev, &stk.ColorF, &stk.ColorB, &stk.Cadence)
		if errScn != nil {
			return &taps.Stakeholder{}, fmt.Errorf("Could not scan stakeholder: %v", errScn)
		}
	}
	if !found {
		return stk, fmt.Errorf("No stakeholder by that email: %w", ErrNotFound)
	}
	return stk, nil
}

func (db *mysqlDB) GetStkDes(email string) (map[string](*taps.Stakeholder), error) {
	topStk, errTop := db.GetStk(email)
	if errTop != nil {
		return map[string](*taps.Stakeholder){}, fmt.Errorf("Could not get top stakeholder %v: %w", email, errTop)
	}
	stks := map[string](*taps.Stakeholder){
		email: topStk,
	}
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	WITH   RECURSIVE des (child, parent) AS
	       (
	       SELECT child
	         ,    parent
	       FROM   stakeholders_hierarchy
	       WHERE  parent = '%v'
	       UNION ALL
	       SELECT h.child
	         ,    h.parent
	       FROM   stakeholders_hierarchy h
	       JOIN   des d
	         ON   h.parent = d.child
	       )
	SELECT s.email
	  ,    s.domain
	  ,    s.name
	  ,    s.abbrev
	  ,    s.colorf
	  ,    s.colorb
	  ,    s.cadence
	FROM   des d
	  JOIN stakeholders s
	  ON   s.email = d.child
	;`, email))
	if errQr != nil {
		return map[string](*taps.Stakeholder){}, fmt.Errorf(
			"Could not query for descendants of stakeholder %v: %v",
			email,
			errQr,
		)
	}
	defer qr.Close()
	for qr.Next() {
		stk := taps.Stakeholder{}
		errScn := qr.Scan(&stk.Email, &stk.Domain, &stk.Name, &stk.Abbrev, &stk.ColorF, &stk.ColorB, &stk.Cadence)
		if errScn != nil {
			return map[string](*taps.Stakeholder){}, fmt.Errorf("Could not scan stakeholder: %v", errScn)
		}
		stks[stk.Email] = &stk
	}
	return stks, nil
}

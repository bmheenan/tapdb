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

func (db *mysqlDB) GetStkAns(email string) (map[string](*taps.Stakeholder), error) {
	btmStk, errBtm := db.GetStk(email)
	if errBtm != nil {
		return map[string](*taps.Stakeholder){}, fmt.Errorf("Could not get bottom stakeholder %v: %w", email, errBtm)
	}
	stks := map[string](*taps.Stakeholder){
		email: btmStk,
	}
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	WITH   RECURSIVE ans (child, parent) AS
	       (
	       SELECT child
	         ,    parent
	       FROM   stakeholders_hierarchy
	       WHERE  child = '%v'
	       UNION ALL
	       SELECT h.child
	         ,    h.parent
	       FROM   stakeholders_hierarchy h
	       JOIN   ans a
	         ON   h.child = a.parent
	       )
	SELECT s.email
	  ,    s.domain
	  ,    s.name
	  ,    s.abbrev
	  ,    s.colorf
	  ,    s.colorb
	  ,    s.cadence
	FROM   ans a
	  JOIN stakeholders s
	  ON   s.email = a.parent
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

func (db *mysqlDB) GetStksForDomain(domain string) (teams []*taps.Team, err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT    s.email
	  ,       s.name
	  ,       s.abbrev
	  ,       s.colorf
	  ,       s.colorb
	  ,       s.cadence
	FROM      stakeholders s
	LEFT JOIN stakeholders_hierarchy h
	  ON      s.email = h.child
	WHERE     h.child IS NULL
	  AND     s.domain = '%v'
	ORDER BY  s.name
	;`, domain))
	if errQr != nil {
		err = fmt.Errorf("Could not query for top level stakeholders: %v", errQr)
		return
	}
	defer qr.Close()
	teams = []*taps.Team{}
	for qr.Next() {
		t := &taps.Team{
			Stk:     taps.Stakeholder{},
			Members: []taps.Team{},
		}
		errScn := qr.Scan(&t.Stk.Email, &t.Stk.Name, &t.Stk.Abbrev, &t.Stk.ColorF, &t.Stk.ColorB, &t.Stk.Cadence)
		if errScn != nil {
			err = fmt.Errorf("Could not scan stakeholder: %v", errScn)
			return
		}
		errCh := db.fillStkChildren(t)
		if errCh != nil {
			err = fmt.Errorf("Could not get children of %v: %v", t.Stk.Email, errCh)
			return
		}
		teams = append(teams, t)
	}
	return
}

func (db *mysqlDB) fillStkChildren(parent *taps.Team) (err error) {
	qr, errQr := db.conn.Query(fmt.Sprintf(`
	SELECT   s.email
	  ,      s.name
	  ,      s.abbrev
	  ,      s.colorf
	  ,      s.colorb
	  ,      s.cadence
	FROM     stakeholders s
	JOIN     stakeholders_hierarchy h
	  ON     s.email = h.child
	WHERE    h.parent = '%v'
	ORDER BY s.name
	;`, parent.Stk.Email))
	if errQr != nil {
		err = fmt.Errorf("Could not query for children of %v: %v", parent.Stk.Email, errQr)
		return
	}
	defer qr.Close()
	parent.Members = []taps.Team{}
	for qr.Next() {
		m := taps.Team{}
		errScn := qr.Scan(&m.Stk.Email, &m.Stk.Name, &m.Stk.Abbrev, &m.Stk.ColorF, &m.Stk.ColorB, &m.Stk.Cadence)
		if errScn != nil {
			err = fmt.Errorf("Could not scan child of %v: %v", parent.Stk.Name, errScn)
			return
		}
		errF := db.fillStkChildren(&m)
		if errF != nil {
			err = fmt.Errorf("Could not get children of %v: %v", m.Stk.Name, errF)
			return
		}
		parent.Members = append(parent.Members, m)
	}
	return
}

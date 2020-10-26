package tapdb

import (
	"fmt"
)

const keyClearDomainPT = "cleardomainpersonteam"
const qryClearDomainPT = `
DELETE FROM
	personteams
WHERE
	domain = ?`
const keyClearDomainPTPC = "cleardomainpersonteamparentchild"
const qryClearDomainPTPC = `
DELETE FROM
	personteams_parent_child
WHERE
	domain = ?`
const keyClearDomainThreads = "cleardomainthreads"
const qryClearDomainThreads = `
DELETE FROM
	threads
WHERE
	domain = ?`
const keyClearDomainThreadsPT = "cleardomainthreadsparentchild"
const qryClearDomainThreadsPT = `
DELETE FROM
	threads_parent_child
WHERE
	domain = ?`

func (db *mySQLDB) initClearDomain() error {
	var err error
	db.stmts[keyClearDomainPT], err = db.conn.Prepare(qryClearDomainPT)
	if err != nil {
		return err
	}
	db.stmts[keyClearDomainPTPC], err = db.conn.Prepare(qryClearDomainPTPC)
	if err != nil {
		return err
	}
	db.stmts[keyClearDomainThreads], err = db.conn.Prepare(qryClearDomainThreads)
	if err != nil {
		return err
	}
	db.stmts[keyClearDomainThreadsPT], err = db.conn.Prepare(qryClearDomainThreadsPT)
	return err
}

func (db *mySQLDB) ClearDomain(dom string) error {
	_, errThreadPC := db.stmts[keyClearDomainThreadsPT].Exec(dom)
	if errThreadPC != nil {
		return fmt.Errorf("Could not delete thread parent/child relationships matching domain %v: %v", dom, errThreadPC)
	}
	_, errThreads := db.stmts[keyClearDomainThreads].Exec(dom)
	if errThreads != nil {
		return fmt.Errorf("Could not delete threads matching domain %v: %v", dom, errThreads)
	}
	_, errPTPC := db.stmts[keyClearDomainPTPC].Exec(dom)
	if errPTPC != nil {
		return fmt.Errorf("Could not delete personteam parent/child relationships matching domain %v: %v", dom, errPTPC)
	}
	_, errPT := db.stmts[keyClearDomainPT].Exec(dom)
	if errPT != nil {
		return fmt.Errorf("Could not delete personteams matching domain %v: %v", dom, errPT)
	}
	return nil
}

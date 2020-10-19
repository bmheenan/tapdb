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

func (db *mySQLDB) initClearDomain() error {
	var err error
	db.stmts[keyClearDomainPT], err = db.conn.Prepare(qryClearDomainPT)
	if err != nil {
		return err
	}
	db.stmts[keyClearDomainPTPC], err = db.conn.Prepare(qryClearDomainPTPC)
	return err
}

func (db *mySQLDB) ClearDomain(dom string) error {
	_, errPT := db.stmts[keyClearDomainPT].Exec(dom)
	if errPT != nil {
		return fmt.Errorf("Could not delete personteams matching domain %v: %v", dom, errPT)
	}
	_, errPTPC := db.stmts[keyClearDomainPTPC].Exec(dom)
	if errPTPC != nil {
		return fmt.Errorf("Could not delete personteam parent/child relationships matching domain %v: %v", dom, errPTPC)
	}
	return nil
}

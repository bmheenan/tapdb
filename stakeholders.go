package tapdb

import (
	"fmt"
)

func (db *mysqlDB) NewStakeholder(thread int64, stakeholder, domain string, ord int, topLvl bool, cost int) error {
	if stakeholder == "" || domain == "" || ord < 0 || cost < 0 {
		return fmt.Errorf("Stakeholder and domain must be non-blank; Ord and cost must be >= 0: %w", ErrBadArgs)
	}
	_, errIn := db.conn.Exec(fmt.Sprintf(`
	INSERT INTO threads_stakeholders
	            (thread, stakeholder, domain, ord, toplevel, costctx)
	VALUES      (    %v,        '%v',   '%v',  %v,       %v,      %v)
	;`, thread, stakeholder, domain, ord, topLvl, cost))
	if errIn != nil {
		return fmt.Errorf("Could not add stakeholder: %v", errIn)
	}
	return nil
}

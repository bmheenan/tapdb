package tapdb

import (
	"errors"
	"fmt"
)

const keyIterationsByPT = "iterationsbypersonteam"
const qryIterationsByPT = `
SELECT
	DISTINCT iteration
  FROM (
    SELECT
        threads.iteration AS iteration,
        threads_stakeholders.stakeholder AS personteam
      FROM
        threads
      INNER JOIN
        threads_stakeholders
      ON
		threads.id = threads_stakeholders.thread
	  WHERE
		threads_stakeholders.stakeholder = ?
    UNION
    SELECT
        iteration,
        owner AS personteam
      FROM
		threads
	  WHERE
		owner = ?
  ) AS all_stakeholders;`

func (db *mySQLDB) initIterationsByPersonteam() error {
	var err error
	db.stmts[keyIterationsByPT], err = db.conn.Prepare(qryIterationsByPT)
	return err
}

func (db *mySQLDB) IterationsByPersonteam(email string) ([]string, error) {
	if email == "" {
		return []string{}, errors.New("Email cannot be blank")
	}
	/*_, errUse := db.conn.Exec(`USE tapestry`)
	if errUse != nil {
		return []string{}, fmt.Errorf("Could not `USE` database: %v", errUse)
	}*/
	result, errQry := db.stmts[keyIterationsByPT].Query(email, email)
	if errQry != nil {
		return []string{}, fmt.Errorf("Could not query for iterations: %v", errQry)
	}
	defer result.Close()
	var iters []string
	for result.Next() {
		var iter string
		errScan := result.Scan(&iter)
		if errScan != nil {
			return []string{}, fmt.Errorf("Could not scan iteration from query result: %v", errScan)
		}
		iters = append(iters, iter)
	}
	return iters, nil
}

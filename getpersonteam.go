package tapdb

import (
	"fmt"

	"github.com/bmheenan/tapstruct"
)

var getPersonteamQuery = `
SELECT
	email,
	name,
	abbrev,
	colorf,
	colorb
FROM
	personteams
WHERE
	email = %v
`

func (t *tapdb) GetPersonteam(email string) (tapstruct.Personteam, error) {
	rows, err := t.db.Query(fmt.Sprintf(getPersonteamQuery, email))
	if err != nil {
		return tapstruct.Personteam{}, fmt.Errorf("Could not query for personteam: %v", err)
	}
	defer rows.Close()
	rows.Next()
	res := tapstruct.Personteam{}
	err = rows.Scan(&res.Email, &res.Name, &res.Abbrev, &res.ColorF, &res.ColorB)
	if err != nil {
		return tapstruct.Personteam{}, fmt.Errorf("Could not get data from query result: %v", err)
	}
	return res, nil
}

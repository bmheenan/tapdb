package tapdb

import (
	"fmt"
)

func (db *mySQLDB) getSubteamMembers(team string) ([]string, error) {
	qrySubteamMembers := `
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
	SELECT DISTINCT d.child
	FROM   descendants d;`
	qrMembers, errQry := db.conn.Query(fmt.Sprintf(qrySubteamMembers, team))
	if errQry != nil {
		return []string{}, fmt.Errorf("Could not query for subteams of %v: %v", team, errQry)
	}
	defer qrMembers.Close()
	members := []string{team}
	for qrMembers.Next() {
		var s string
		errScn := qrMembers.Scan(&s)
		if errScn != nil {
			return []string{}, fmt.Errorf("Could not scan sub team member: %v", errScn)
		}
		members = append(members, s)
	}
	return members, nil
}

func (db *mySQLDB) concatStringAsList(iters []string) string {
	res := ""
	l := len(iters) - 1
	for i, v := range iters {
		res = res + fmt.Sprintf("'%s'", v)
		if i < l {
			res = res + ", "
		}
	}
	return res
}

func (db *mySQLDB) concatInt64AsList(iters []int64) string {
	res := ""
	l := len(iters) - 1
	for i, v := range iters {
		res = res + fmt.Sprintf("'%d'", v)
		if i < l {
			res = res + ", "
		}
	}
	return res
}

func (db *mySQLDB) min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (db *mySQLDB) max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (db *mySQLDB) abs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}

func (db *mySQLDB) strIncl(strs *[]string, s string) bool {
	for _, str := range *strs {
		if str == s {
			return true
		}
	}
	return false
}

/*
func (db *mySQLDB) ptIncludes(pt *tapstruct.Personteam, s string) bool {
	if pt.Email == s {
		return true
	}
	for _, c := range pt.Children {
		if db.ptIncludes(&c, s) {
			return true
		}
	}
	return false
}


func (db *mySQLDB) removeInt(a *[]int64, item int64) {
	index := -1
	for i, v := range *a {
		if v == item {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	copy((*a)[index:], (*a)[index+1:])
	(*a)[len(*a)-1] = 0
	*a = (*a)[:len(*a)-1]
}


func (db *mySQLDB) removeThread(a *[]*threadWMeta, i int) {
	copy((*a)[i:], (*a)[i+1:])
	(*a)[len(*a)-1] = nil
	*a = (*a)[:len(*a)-1]
}
*/

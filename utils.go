package tapdb

import "fmt"

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
		res = res + fmt.Sprintf("'%s'", v)
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

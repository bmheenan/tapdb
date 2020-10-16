package tapdb

import "fmt"

type connVars struct {
	host       string
	port       string
	unixSocket string
	user       string
	pass       string
	dbName     string
}

func (cv *connVars) formatName() string {
	var cred string
	if cv.user != "" {
		cred = cv.user
		if cv.pass != "" {
			cred = cred + ":" + cv.pass
		}
		cred = cred + "@"
	}
	if cv.unixSocket != "" {
		return fmt.Sprintf("%sunix(%s)/%s", cred, cv.unixSocket, cv.dbName)
	}
	return fmt.Sprintf("%stcp([%s]:%s)/%s", cred, cv.host, cv.port, cv.dbName)
}

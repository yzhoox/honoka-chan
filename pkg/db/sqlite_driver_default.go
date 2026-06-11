//go:build !android

package db

import (
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

const sqliteDriverName = "sqlite"

func mainSQLiteDSN(dbPath string) string {
	return dbPath
}

func userSQLiteDSN(dbPath string) string {
	return dbPath +
		"?_pragma=busy_timeout(5000)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=synchronous(NORMAL)"
}

func prepareMainSQLiteEngine(*xorm.Engine) error {
	return nil
}

func prepareUserSQLiteEngine(*xorm.Engine) error {
	return nil
}

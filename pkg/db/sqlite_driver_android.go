//go:build android

package db

import (
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

const sqliteDriverName = "sqlite3"

func mainSQLiteDSN(dbPath string) string {
	return dbPath
}

func userSQLiteDSN(dbPath string) string {
	return dbPath
}

func prepareMainSQLiteEngine(engine *xorm.Engine) error {
	return applySQLitePragmas(engine, nil)
}

func prepareUserSQLiteEngine(engine *xorm.Engine) error {
	return applySQLitePragmas(engine, []string{
		"PRAGMA busy_timeout = 5000",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
	})
}

func applySQLitePragmas(engine *xorm.Engine, pragmas []string) error {
	for _, pragma := range pragmas {
		if _, err := engine.Exec(pragma); err != nil {
			return err
		}
	}
	return nil
}

package db

import (
	"sync"
	"xorm.io/xorm"
)

var (
	MainDb  = "assets/main.db"
	UserDb  = "assets/data.db"
	MainEng *xorm.Engine
	UserEng *xorm.Engine

	mu sync.Mutex
)

func Init(mainDbPath, userDbPath string) error {
	mu.Lock()
	defer mu.Unlock()

	if mainDbPath == "" {
		mainDbPath = MainDb
	}
	if userDbPath == "" {
		userDbPath = UserDb
	}

	if MainEng != nil || UserEng != nil {
		if MainDb == mainDbPath && UserDb == userDbPath {
			return nil
		}
		if err := closeLocked(); err != nil {
			return err
		}
	}

	mainEng, err := xorm.NewEngine(sqliteDriverName, mainSQLiteDSN(mainDbPath))
	if err != nil {
		return err
	}
	if err := prepareMainSQLiteEngine(mainEng); err != nil {
		mainEng.Close()
		return err
	}
	if err := mainEng.Ping(); err != nil {
		mainEng.Close()
		return err
	}
	mainEng.SetMaxOpenConns(10)
	mainEng.SetMaxIdleConns(5)

	userEng, err := xorm.NewEngine(sqliteDriverName, userSQLiteDSN(userDbPath))
	if err != nil {
		mainEng.Close()
		return err
	}
	if err := prepareUserSQLiteEngine(userEng); err != nil {
		mainEng.Close()
		userEng.Close()
		return err
	}
	if err := userEng.Ping(); err != nil {
		mainEng.Close()
		userEng.Close()
		return err
	}
	userEng.SetMaxOpenConns(10)
	userEng.SetMaxIdleConns(5)

	MainDb = mainDbPath
	UserDb = userDbPath
	MainEng = mainEng
	UserEng = userEng
	return nil
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()
	return closeLocked()
}

func closeLocked() error {
	var firstErr error

	if MainEng != nil {
		if err := MainEng.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		MainEng = nil
	}
	if UserEng != nil {
		if err := UserEng.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		UserEng = nil
	}

	return firstErr
}

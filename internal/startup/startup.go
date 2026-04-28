package startup

func StartUp() {
	CreateTables()
	LoadUnitData()
	CreateDefaultUser()
}

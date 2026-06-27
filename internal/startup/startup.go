package startup

func StartUp() {
	CreateTables()
	MigrateLegacyUnitTables()
	LoadUnitData()
	CreateDefaultUser()
	EnsureDefaultFriends()
}

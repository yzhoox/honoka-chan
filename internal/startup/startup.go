package startup

func StartUp() {
	CreateTables()
	LoadUnitData()
	MigrateLegacyUnitTables()
	CreateDefaultUser()
	EnsureDefaultFriends()
}

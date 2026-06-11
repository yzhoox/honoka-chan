package startup

func StartUp() error {
	if err := CreateTables(); err != nil {
		return err
	}
	if err := LoadUnitData(); err != nil {
		return err
	}
	if err := CreateDefaultUser(); err != nil {
		return err
	}
	if err := EnsureDefaultFriends(); err != nil {
		return err
	}
	return nil
}

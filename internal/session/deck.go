package session

import (
	usermodel "honoka-chan/internal/model/user"
)

func (ss *Session) GetUserDeck(deckID int) (bool, *usermodel.UserDeck) {
	deckData := usermodel.UserDeck{}
	has, err := ss.UserEng.Table("user_deck").
		Where("user_id = ? AND deck_id = ?", ss.UserID, deckID).Get(&deckData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &deckData
}

func (ss *Session) GetUserDeckUnit(deckID int) []usermodel.UserDeckUnit {
	has, deckData := ss.GetUserDeck(deckID)
	if !has {
		return []usermodel.UserDeckUnit{}
	}

	unitData := []usermodel.UserDeckUnit{}
	err := ss.UserEng.Table("user_deck_unit").
		Where("user_deck_id = ?", deckData.ID).Find(&unitData)
	if ss.CheckErr(err) {
		return []usermodel.UserDeckUnit{}
	}

	return unitData
}

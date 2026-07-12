package session

import (
	usermodel "honoka-chan/internal/model/user"
)

func (ss *Session) GetUserDeck(deckID int) (bool, *usermodel.UserDeck, error) {
	deckData := usermodel.UserDeck{}
	has, err := ss.UserEng.Table("user_deck").
		Where("user_id = ? AND deck_id = ?", ss.UserID, deckID).Get(&deckData)
	if err != nil {
		return false, nil, err
	}

	return has, &deckData, nil
}

func (ss *Session) GetUserDeckUnit(deckID int) []usermodel.UserDeckUnit {
	has, deckData, err := ss.GetUserDeck(deckID)
	if ss.CheckErr(err) {
		return []usermodel.UserDeckUnit{}
	}
	if !has {
		return []usermodel.UserDeckUnit{}
	}

	unitData := []usermodel.UserDeckUnit{}
	err = ss.UserEng.Table("user_deck_unit").
		Where("user_deck_id = ? AND user_id = ?", deckData.ID, ss.UserID).Find(&unitData)
	if ss.CheckErr(err) {
		return []usermodel.UserDeckUnit{}
	}

	return unitData
}

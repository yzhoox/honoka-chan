package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func unitDeckInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	userDeck := []unitapischema.UserDeckData{}
	err = ss.UserEng.Table("user_deck").Where("user_id = ?", ss.UserID).Asc("deck_id").Find(&userDeck)
	if ss.CheckErr(err) {
		return
	}

	unitDeckInfo := []unitapischema.DeckInfoData{}
	for _, deck := range userDeck {
		deckUnit := []unitapischema.UnitDeckData{}
		err = ss.UserEng.Table("user_deck_unit").Where("user_deck_id = ?", deck.ID).Asc("position").Find(&deckUnit)
		if ss.CheckErr(err) {
			return
		}

		oUID := []unitapischema.UnitOwningUserIds{}
		for _, u := range deckUnit {
			oUID = append(oUID, unitapischema.UnitOwningUserIds{
				Position:         u.Position,
				UnitOwningUserID: u.UnitOwningUserID,
			})
		}

		mainFlag := false
		if deck.MainFlag == 1 {
			mainFlag = true
		}
		unitDeckInfo = append(unitDeckInfo, unitapischema.DeckInfoData{
			UnitDeckID:        deck.DeckID,
			MainFlag:          mainFlag,
			DeckName:          deck.DeckName,
			UnitOwningUserIds: oUID,
		})
	}
	res = unitapischema.DeckInfoResp{
		Result:     unitDeckInfo,
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}

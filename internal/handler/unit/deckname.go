package unit

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func deckName(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	deckReq := unit.DeckNameReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &deckReq)
	if ss.CheckErr(err) {
		return
	}

	_, err = ss.UserEng.Table("user_deck").Where("user_id = ? AND deck_id = ?", ss.UserID, deckReq.UnitDeckID).Exist()
	if ss.CheckErr(err) {
		return
	}

	userDeck := unitapischema.UserDeckData{
		DeckName: deckReq.DeckName,
	}
	_, err = ss.UserEng.Table("user_deck").Update(&userDeck, &unitapischema.UserDeckData{
		UserID: ss.UserID,
		DeckID: deckReq.UnitDeckID,
	})
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(unit.DeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/deckName", middleware.Common, deckName)
}

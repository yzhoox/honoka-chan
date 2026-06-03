package unit

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	unitapischema "honoka-chan/internal/schema/api/unit"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func deckName(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	deckReq := unitschema.DeckNameReq{}
	err := honokautils.ParseRequestData(ctx, &deckReq)
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

	ss.Respond(unitschema.DeckResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/deckName", middleware.Common, deckName)
}

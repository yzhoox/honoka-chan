package live

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func partyList(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()
	userInfo := ss.GetUserInfo()

	pref := usermodel.UserPref{}
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if ss.CheckErr(err) {
		return
	}

	supportRows, err := listLiveSupportRows(ss)
	if ss.CheckErr(err) {
		return
	}

	partyList := make([]liveschema.PartyList, 0, len(supportRows))
	for _, row := range supportRows {
		party, err := buildLiveSupportParty(ss, row)
		if ss.CheckErr(err) {
			return
		}
		partyList = append(partyList, party)
	}

	if len(partyList) == 0 {
		selfParty := liveschema.PartyList{
			UserInfo:       userInfo,
			SettingAwardID: pref.AwardID,
			FriendStatus:   usermodel.ClientFriendStatusNone,
		}
		selfParty.AvailableSocialPoint = 10
		selfParty.CenterUnitInfo, err = mustFindCenterUnitInfo(ss, ss.UserID, pref.UnitOwningUserID)
		if ss.CheckErr(err) {
			return
		}
		partyList = append(partyList, selfParty)
	}

	ss.Respond(liveschema.PartyListResp{
		ResponseData: liveschema.PartyListData{
			PartyList:         partyList,
			TrainingEnergy:    userInfo.TrainingEnergy,
			TrainingEnergyMax: userInfo.TrainingEnergyMax,
			ServerTimestamp:   time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/live/partyList", middleware.Common, partyList)
}

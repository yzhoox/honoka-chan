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
	partyListData, err := BuildPartyListData(ss)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(liveschema.PartyListResp{
		ResponseData: partyListData,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func BuildPartyListData(ss *session.Session) (liveschema.PartyListData, error) {
	userInfo := ss.GetUserInfo()

	pref := usermodel.UserPref{}
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if err != nil {
		return liveschema.PartyListData{}, err
	}

	supportRows, err := listLiveSupportRows(ss)
	if err != nil {
		return liveschema.PartyListData{}, err
	}

	partyList := make([]liveschema.PartyList, 0, len(supportRows))
	for _, row := range supportRows {
		party, err := buildLiveSupportParty(ss, row)
		if err != nil {
			return liveschema.PartyListData{}, err
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
		if err != nil {
			return liveschema.PartyListData{}, err
		}
		partyList = append(partyList, selfParty)
	}

	return liveschema.PartyListData{
		PartyList:         partyList,
		TrainingEnergy:    userInfo.TrainingEnergy,
		TrainingEnergyMax: userInfo.TrainingEnergyMax,
		ServerTimestamp:   time.Now().Unix(),
	}, nil
}

func init() {
	router.AddHandler("main.php", "POST", "/live/partyList", middleware.Common, partyList)
}

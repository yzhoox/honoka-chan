package handler

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"honoka-chan/internal/tools"
	"honoka-chan/internal/utils"
	honokautils "honoka-chan/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EventRes struct {
	EventScenarioId int    `xorm:"event_scenario_id"`
	Chapter         int    `xorm:"chapter"`
	ChapterAsset    string `xorm:"chapter_asset"`
	OpenDate        string `xorm:"open_date"`
}

type MultiRes struct {
	MultiUnitScenarioId       int    `xorm:"multi_unit_scenario_id"`
	Chapter                   int    `xorm:"chapter"`
	MultiUnitScenarioBtnAsset string `xorm:"multi_unit_scenario_btn_asset"`
	OpenDate                  string `xorm:"open_date"`
}

type MuseumRes struct {
	MuseumContentsId int `xorm:"museum_contents_id"`
	SmileBuff        int `xorm:"smile_buff"`
	PureBuff         int `xorm:"pure_buff"`
	CoolBuff         int `xorm:"cool_buff"`
}

func Api(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	apiReq := []model.ApiReq{}
	err := json.Unmarshal([]byte(ctx.GetString("request_data")), &apiReq)
	if ss.CheckErr(err) {
		return
	}

	results := []any{}
	for _, v := range apiReq {
		// fmt.Println(v)
		// fmt.Println(v.Module, v.Action)

		switch v.Module {
		case "login":
			switch v.Action {
			case "topInfo":
				// key = "login_topinfo_result"
				topInfoResp := model.TopInfoResp{
					Result: model.TopInfoRes{
						FriendActionCnt:        0,
						FriendGreetCnt:         0,
						FriendVarietyCnt:       0,
						FriendNewCnt:           0,
						PresentCnt:             0,
						SecretBoxBadgeFlag:     false,
						ServerDatetime:         time.Now().Format("2006-01-02 15:04:05"),
						ServerTimestamp:        time.Now().Unix(),
						NoticeFriendDatetime:   time.Now().Format("2006-01-02 15:04:05"),
						NoticeMailDatetime:     "2000-01-01 12:00:00",
						FriendsApprovalWaitCnt: 0,
						FriendsRequestCnt:      0,
						IsTodayBirthday:        false,
						LicenseInfo: model.TopInfoLicenseInfo{
							LicenseList:  []any{},
							LicensedInfo: []any{},
							ExpiredInfo:  []any{},
							BadgeFlag:    false,
						},
						UsingBuffInfo:     []any{},
						IsKlabIDTaskFlag:  false,
						KlabIDTaskCanSync: false,
						HasUnreadAnnounce: false,
						ExchangeBadgeCnt:  []int{0, 0, 0},
						AdFlag:            false,
						HasAdReward:       false,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, topInfoResp)
			case "topInfoOnce":
				// key = "login_topinfo_once_result"
				topInfoOnceResp := model.TopInfoOnceResp{
					Result: model.TopInfoOnceRes{
						NewAchievementCnt:            0,
						UnaccomplishedAchievementCnt: 0,
						LiveDailyRewardExist:         false,
						TrainingEnergy:               10,
						TrainingEnergyMax:            10,
						Notification: model.TopInfoOnceNotification{
							Push:       false,
							Lp:         false,
							UpdateInfo: false,
							Campaign:   false,
							Live:       false,
							Lbonus:     false,
							Event:      false,
							Secretbox:  false,
							Birthday:   true,
						},
						OpenArena:               true,
						CostumeStatus:           true,
						OpenAccessory:           true,
						ArenaSiSkillUniqueCheck: true,
						OpenV98:                 true,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, topInfoOnceResp)
			default:
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "live":
			switch v.Action {
			case "liveStatus":
				// key = "live_status_result"
				var liveDifficultyId []int
				normalLives := []model.NormalLiveStatusList{}
				err = ss.MainEng.Table("normal_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyId)
				if ss.CheckErr(err) {
					return
				}
				for _, id := range liveDifficultyId {
					normalLive := model.NormalLiveStatusList{
						LiveDifficultyID:   id,
						Status:             1,
						HiScore:            0,
						HiComboCount:       0,
						ClearCnt:           0,
						AchievedGoalIDList: []int{},
					}
					normalLives = append(normalLives, normalLive)
				}

				specialLives := []model.SpecialLiveStatusList{}
				err = ss.MainEng.Table("special_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyId)
				if ss.CheckErr(err) {
					return
				}
				for _, id := range liveDifficultyId {
					specialLive := model.SpecialLiveStatusList{
						LiveDifficultyID:   id,
						Status:             1,
						HiScore:            0,
						HiComboCount:       0,
						ClearCnt:           0,
						AchievedGoalIDList: []int{},
					}
					specialLives = append(specialLives, specialLive)
				}

				LiveStatusResp := model.LiveStatusResp{
					Result: model.LiveStatusRes{
						NormalLiveStatusList:   normalLives,
						SpecialLiveStatusList:  specialLives,
						TrainingLiveStatusList: []model.TrainingLiveStatusList{},
						MarathonLiveStatusList: []any{},
						FreeLiveStatusList:     []any{},
						CanResumeLive:          false,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, LiveStatusResp)
			case "schedule":
				// key = "live_list_result"
				var liveDifficultyId []int
				specialLives := []model.SpecialLiveStatusList{}
				err = ss.MainEng.Table("special_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyId)
				if ss.CheckErr(err) {
					return
				}
				for _, id := range liveDifficultyId {
					specialLive := model.SpecialLiveStatusList{
						LiveDifficultyID:   id,
						Status:             1,
						HiScore:            0,
						HiComboCount:       0,
						ClearCnt:           0,
						AchievedGoalIDList: []int{},
					}
					specialLives = append(specialLives, specialLive)
				}

				livesList := []model.LiveList{}
				for _, v := range specialLives {
					livesList = append(livesList, model.LiveList{
						LiveDifficultyID: v.LiveDifficultyID,
						StartDate:        "2023-01-01 00:00:00",
						EndDate:          "2037-01-01 00:00:00",
						IsRandom:         false,
					})
				}
				liveListResp := model.LiveScheduleResp{
					Result: model.LiveScheduleRes{
						EventList:              []any{},
						LiveList:               livesList,
						LimitedBonusList:       []any{},
						LimitedBonusCommonList: []model.LimitedBonusCommonList{}, // 特效道具
						RandomLiveList:         []model.RandomLiveList{},         // 随机歌曲
						FreeLiveList:           []any{},
						TrainingLiveList:       []model.TrainingLiveList{}, // 挑战歌曲
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, liveListResp)
			default:
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "unit":
			switch v.Action {
			case "unitAll":
				// key = "unit_list_result"
				unitsData := []model.Active{}
				err = ss.MainEng.Table("common_unit_m").Find(&unitsData)
				if ss.CheckErr(err) {
					return
				}

				userUnits := []model.Active{}
				err = ss.UserEng.Table("user_unit_m").Where("user_id = ?", ctx.GetString("userid")).Find(&userUnits)
				if ss.CheckErr(err) {
					return
				}
				unitsData = append(unitsData, userUnits...)

				unitListResp := model.UnitAllResp{
					Result: model.UnitAllRes{
						Active:  unitsData,
						Waiting: []model.Waiting{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, unitListResp)
			case "deckInfo":
				// key = "unit_deck_result"
				userDeck := []model.UserDeckData{}
				err = ss.UserEng.Table("user_deck_m").Where("user_id = ?", ctx.GetString("userid")).Asc("deck_id").Find(&userDeck)
				if ss.CheckErr(err) {
					return
				}

				unitDeckInfo := []model.UnitDeckInfoRes{}
				for _, deck := range userDeck {
					deckUnit := []model.UnitDeckData{}
					err = ss.UserEng.Table("deck_unit_m").Where("user_deck_id = ?", deck.ID).Asc("position").Find(&deckUnit)
					if ss.CheckErr(err) {
						return
					}

					oUids := []model.UnitOwningUserIds{}
					for _, unit := range deckUnit {
						oUids = append(oUids, model.UnitOwningUserIds{
							Position:         unit.Position,
							UnitOwningUserID: unit.UnitOwningUserID,
						})
					}

					mainFlag := false
					if deck.MainFlag == 1 {
						mainFlag = true
					}
					unitDeckInfo = append(unitDeckInfo, model.UnitDeckInfoRes{
						UnitDeckID:        deck.DeckID,
						MainFlag:          mainFlag,
						DeckName:          deck.DeckName,
						UnitOwningUserIds: oUids,
					})
				}
				unitDeckResp := model.UnitDeckInfoResp{
					Result:     unitDeckInfo,
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, unitDeckResp)
			case "supporterAll":
				// key = "unit_support_result"
				unitSupportResp := model.UnitSupportResp{
					Result: model.UnitSupportRes{
						UnitSupportList: []model.UnitSupportList{},
					}, // 练习道具
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, unitSupportResp)
			case "removableSkillInfo":
				// key = "owning_equip_result"
				var skillEquipCount []model.SkillEquipCount
				err := ss.UserEng.Table("skill_equip_m").Where("user_id = ?", ctx.GetString("userid")).Select("unit_removable_skill_id,COUNT(*) AS ct").
					GroupBy("unit_removable_skill_id").Find(&skillEquipCount)
				if ss.CheckErr(err) {
					return
				}

				var rmSkillIds []int
				err = ss.MainEng.Table("unit_removable_skill_m").Where("effect_range = 1").Cols("unit_removable_skill_id").Find(&rmSkillIds)
				if ss.CheckErr(err) {
					return
				}

				owingInfo := []model.OwningInfo{}
				for _, id := range rmSkillIds {
					info := model.OwningInfo{
						UnitRemovableSkillID: id,
						TotalAmount:          9,
						EquippedAmount:       0,
						InsertDate:           "2023-01-01 12:00:00",
					}
					for _, sk := range skillEquipCount {
						if id == sk.UnitRemovableSkillId {
							info.EquippedAmount = sk.Count
							break
						}
					}
					owingInfo = append(owingInfo, info)
				}

				var unitOwningIds []int
				err = ss.UserEng.Table("skill_equip_m").Where("user_id = ?", ctx.GetString("userid")).Cols("unit_owning_user_id").GroupBy("unit_owning_user_id").Find(&unitOwningIds)
				if ss.CheckErr(err) {
					return
				}

				equipInfo := map[int]any{}
				for _, v := range unitOwningIds {
					detail := []model.SkillEquipDetail{}
					err = ss.UserEng.Table("skill_equip_m").Where("user_id = ? AND unit_owning_user_id = ?", ctx.GetString("userid"), v).
						Cols("unit_removable_skill_id").Find(&detail)
					if ss.CheckErr(err) {
						return
					}

					equipInfo[v] = model.SkillEquipList{
						UnitOwningUserID: v,
						Detail:           detail,
					}
				}

				rmSkillResp := model.RemovableSkillResp{
					Result: model.RemovableSkillRes{
						OwningInfo:    owingInfo,
						EquipmentInfo: equipInfo,
					}, // 宝石
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, rmSkillResp)
			case "accessoryAll":
				// key = "unit_accessory_result"
				accessoryList := []model.AccessoryList{}
				err := ss.MainEng.Table("common_accessory_m").Find(&accessoryList)
				if ss.CheckErr(err) {
					return
				}
				for k := range accessoryList {
					accessoryList[k].NextExp = 0
					accessoryList[k].Level = 8
					accessoryList[k].MaxLevel = 8
					accessoryList[k].RankUpCount = 4
					accessoryList[k].FavoriteFlag = true
				}
				wearingInfo := []model.WearingInfo{}
				err = ss.UserEng.Table("accessory_wear_m").Where("user_id = ?", ctx.GetString("userid")).Find(&wearingInfo)
				if ss.CheckErr(err) {
					return
				}
				unitAccResp := model.UnitAccessoryAllResp{
					Result: model.UnitAccessoryAllResult{
						AccessoryList:      accessoryList,
						WearingInfo:        wearingInfo,
						EspecialCreateFlag: false,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, unitAccResp)
			default:
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "costume":
			if v.Action == "costumeList" {
				// key = "costume_list_result"
				costumeListResp := model.CostumeListResp{
					Result: model.CostumeListRes{
						CostumeList: []model.CostumeList{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, costumeListResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "album":
			if v.Action == "albumAll" {
				// key = "album_unit_result"
				albumLists := []model.AlbumResult{}
				unitList := []AlbumRes{}
				err = ss.MainEng.Table("unit_m").Cols("unit_id,rarity").OrderBy("unit_id ASC").Find(&unitList)
				if ss.CheckErr(err) {
					return
				}

				for _, unit := range unitList {
					albumList := model.AlbumResult{
						RankMaxFlag:      true,
						LoveMaxFlag:      true,
						RankLevelMaxFlag: true,
						AllMaxFlag:       true,
						FavoritePoint:    1000,
					}
					albumList.UnitID = unit.UnitId
					if unit.Rarity != 4 {
						albumList.SignFlag = false
						switch unit.Rarity {
						case 1:
							albumList.HighestLovePerUnit = 50
							albumList.TotalLove = 50
						case 2:
							albumList.HighestLovePerUnit = 200
							albumList.TotalLove = 200
						case 3:
							albumList.HighestLovePerUnit = 500
							albumList.TotalLove = 500
						case 5:
							albumList.HighestLovePerUnit = 750
							albumList.TotalLove = 750
						}
					} else {
						albumList.HighestLovePerUnit = 1000
						albumList.TotalLove = 1000

						// IsSigned
						albumList.SignFlag = utils.IsSigned(unit.UnitId)
					}
					albumLists = append(albumLists, albumList)
				}

				albumResp := model.AlbumResp{
					Result:     albumLists,
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, albumResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "scenario":
			if v.Action == "scenarioStatus" {
				// key = "scenario_status_result"
				var scenarioIds []int
				scenarioLists := []model.ScenarioStatusList{}
				err = ss.MainEng.Table("scenario_m").Cols("scenario_id").OrderBy("scenario_id ASC").Find(&scenarioIds)
				if ss.CheckErr(err) {
					return
				}

				for _, id := range scenarioIds {
					scenarioLists = append(scenarioLists, model.ScenarioStatusList{
						ScenarioID: id,
						Status:     2,
					})
				}
				scenarioResp := model.ScenarioStatusResp{
					Result: model.ScenarioStatusRes{
						ScenarioStatusList: scenarioLists,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, scenarioResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "subscenario":
			if v.Action == "subscenarioStatus" {
				// key = "subscenario_status_result"
				var subScenarioIds []int
				subScenarioLists := []model.SubscenarioStatusList{}
				err = ss.MainEng.Table("subscenario_m").Cols("subscenario_id").OrderBy("subscenario_id ASC").Find(&subScenarioIds)
				if ss.CheckErr(err) {
					return
				}

				for _, id := range subScenarioIds {
					subScenarioLists = append(subScenarioLists, model.SubscenarioStatusList{
						SubscenarioID: id,
						Status:        2,
					})
				}
				subScenarioResp := model.SubscenarioStatusResp{
					Result: model.SubscenarioStatusRes{
						SubscenarioStatusList:  subScenarioLists,
						UnlockedSubscenarioIds: []any{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, subScenarioResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "eventscenario":
			if v.Action == "status" {
				// key = "event_scenario_result"
				var eventIds []int
				eventsList := []model.EventScenarioList{}
				err = ss.MainEng.Table("event_scenario_m").Cols("event_id").GroupBy("event_id").OrderBy("event_id DESC").Find(&eventIds)
				if ss.CheckErr(err) {
					return
				}

				for _, id := range eventIds {
					eventRes := []EventRes{}
					chapsList := []model.EventScenarioChapterList{}
					err = ss.MainEng.Table("event_scenario_m").Where("event_id = ?", id).Cols("event_scenario_id,chapter,chapter_asset,open_date").
						OrderBy("chapter DESC").Find(&eventRes)
					if ss.CheckErr(err) {
						return
					}

					for _, res := range eventRes {
						chapList := model.EventScenarioChapterList{
							EventScenarioID: res.EventScenarioId,
							Chapter:         res.Chapter,
							ChapterAsset:    res.ChapterAsset,
							Status:          2,
							OpenFlashFlag:   0,
							IsReward:        false,
							CostType:        1000,
							ItemID:          1200,
							Amount:          1,
						}
						chapsList = append(chapsList, chapList)
					}

					eventList := model.EventScenarioList{
						EventID:     id,
						OpenDate:    strings.ReplaceAll(eventRes[0].OpenDate, "/", "-"),
						ChapterList: chapsList,
					}

					// HACK event_scenario_btn_asset
					switch id {
					case 10001:
						eventList.EventScenarioBtnAsset = "assets/image/ui/eventscenario/38_se_ba_t.png"
					case 221:
						eventList.EventScenarioBtnAsset = "assets/image/ui/eventscenario/215_se_ba_t.png"
					default:
						eventList.EventScenarioBtnAsset = fmt.Sprintf("assets/image/ui/eventscenario/%d_se_ba_t.png", id)
					}

					eventsList = append(eventsList, eventList)
				}
				eventScenarioResp := model.EventScenarioStatusResp{
					Result: model.EventScenarioStatusRes{
						EventScenarioList: eventsList,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, eventScenarioResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "multiunit":
			if v.Action == "multiunitscenarioStatus" {
				// key = "multi_unit_scenario_result"
				var multiIds []int
				multiUnitsList := []model.MultiUnitScenarioStatusList{}
				err = ss.MainEng.Table("multi_unit_scenario_m").Cols("multi_unit_id").GroupBy("multi_unit_id").OrderBy("multi_unit_id ASC").Find(&multiIds)
				if ss.CheckErr(err) {
					return
				}

				for _, id := range multiIds {
					multiRes := MultiRes{}
					_, err = ss.MainEng.Table("multi_unit_scenario_m").
						Join("LEFT", "multi_unit_scenario_open_m", "multi_unit_scenario_m.multi_unit_id = multi_unit_scenario_open_m.multi_unit_id").
						Cols("multi_unit_scenario_btn_asset,open_date,multi_unit_scenario_id,chapter").
						Where("multi_unit_scenario_m.multi_unit_id = ?", id).Get(&multiRes)
					if ss.CheckErr(err) {
						return
					}

					multiUnitsList = append(multiUnitsList, model.MultiUnitScenarioStatusList{
						MultiUnitID:               id,
						Status:                    2,
						MultiUnitScenarioBtnAsset: multiRes.MultiUnitScenarioBtnAsset,
						OpenDate:                  strings.ReplaceAll(multiRes.OpenDate, "/", "-"),
						ChapterList: []model.MultiUnitScenarioChapterList{
							{
								MultiUnitScenarioID: multiRes.MultiUnitScenarioId,
								Chapter:             multiRes.Chapter,
								Status:              2,
							},
						},
					})
				}
				unitsResp := model.MultiUnitScenarioStatusResp{
					Result: model.MultiUnitScenarioStatusRes{
						MultiUnitScenarioStatusList:  multiUnitsList,
						UnlockedMultiUnitScenarioIds: []any{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, unitsResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "payment":
			if v.Action == "productList" {
				// key = "product_result"
				productResp := model.ProductListResp{
					Result: model.ProductListRes{
						RestrictionInfo: model.RestrictionInfo{
							Restricted: false,
						},
						UnderAgeInfo: model.UnderAgeInfo{
							BirthSet:    false,
							HasLimit:    false,
							LimitAmount: nil,
							MonthUsed:   0,
						},
						SnsProductList:   []model.SnsProductList{},
						ProductList:      []model.ProductList{},
						SubscriptionList: []model.SubscriptionList{},
						ShowPointShop:    false,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, productResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "banner":
			if v.Action == "bannerList" {
				// key = "banner_result"
				bannerResp := model.BannerListResp{
					Result: model.BannerListRes{
						TimeLimit: "2037-12-31 23:59:59",
						BannerList: []model.BannerList{
							{
								BannerType:       1,
								TargetID:         1743,
								AssetPath:        "assets/image/secretbox/icon/s_ba_1718_1.png",
								FixedFlag:        false,
								BackSide:         false,
								BannerID:         101151,
								StartDate:        "2013-04-15 00:00:00",
								EndDate:          "2037-12-31 23:59:59",
								AddUnitStartDate: "2022-01-01 00:00:00",
							},
							{
								BannerType:       1,
								TargetID:         1741,
								AssetPath:        "assets/image/secretbox/icon/s_ba_1719_1.png",
								FixedFlag:        false,
								BackSide:         false,
								BannerID:         101150,
								StartDate:        "2013-04-15 00:00:00",
								EndDate:          "2037-12-31 23:59:59",
								AddUnitStartDate: "2022-01-01 00:00:00",
							},
							{
								BannerType:       1,
								TargetID:         1740,
								AssetPath:        "assets/image/secretbox/icon/s_ba_1720_1.png",
								FixedFlag:        false,
								BackSide:         false,
								BannerID:         101149,
								StartDate:        "2013-04-15 00:00:00",
								EndDate:          "2037-12-31 23:59:59",
								AddUnitStartDate: "2022-01-01 00:00:00",
							},
							{
								BannerType:       1,
								TargetID:         1739,
								AssetPath:        "assets/image/secretbox/icon/s_ba_1721_1.png",
								FixedFlag:        false,
								BackSide:         false,
								BannerID:         101144,
								StartDate:        "2013-04-15 00:00:00",
								EndDate:          "2037-12-31 23:59:59",
								AddUnitStartDate: "2022-01-01 00:00:00",
							},
							{
								BannerType: 2,
								TargetID:   1,
								AssetPath:  "assets/image/webview/wv_ba_01.png",
								WebviewURL: "/manga",
								FixedFlag:  false,
								BackSide:   true,
								BannerID:   200001,
								StartDate:  "2016-10-15 15:00:00",
								EndDate:    "2037-12-31 23:59:59",
							},
						},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, bannerResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "notice":
			if v.Action == "noticeMarquee" {
				// key = "item_marquee_result"
				marqueeResp := model.NoticeMarqueeResp{
					Result: model.NoticeMarqueeRes{
						ItemCount:   0,
						MarqueeList: []any{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, marqueeResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "user":
			switch v.Action {
			case "getNavi":
				// key = "user_intro_result"
				var uId, oId int
				_, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Cols("user_id,unit_owning_user_id").Get(&uId, &oId)
				if ss.CheckErr(err) {
					return
				}

				userIntroResp := model.UserNaviResp{
					Result: model.UserNaviRes{
						User: model.User{
							UserID:           uId,
							UnitOwningUserID: oId,
						},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, userIntroResp)
			case "userInfo":
				userId, err := strconv.Atoi(ctx.GetString("userid"))
				if ss.CheckErr(err) {
					return
				}

				pref := tools.UserPref{}
				exists, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", userId).Get(&pref)
				if ss.CheckErr(err) {
					return
				}
				if !exists {
					ctx.String(http.StatusForbidden, ErrorMsg)
					return
				}

				userInfoResp := model.ApiUserInfoResp{
					Result: model.UserInfo{
						UserID:                         userId,
						Name:                           pref.UserName,
						Level:                          config.Conf.UserPrefs.Level,
						Exp:                            config.Conf.UserPrefs.ExpNumerator,
						PreviousExp:                    0,
						NextExp:                        config.Conf.UserPrefs.ExpDenominator,
						GameCoin:                       config.Conf.UserPrefs.GameCoin,
						SnsCoin:                        config.Conf.UserPrefs.SnsCoin,
						FreeSnsCoin:                    config.Conf.UserPrefs.SnsCoin,
						PaidSnsCoin:                    0,
						SocialPoint:                    1438395,
						UnitMax:                        5000,
						WaitingUnitMax:                 1000,
						EnergyMax:                      config.Conf.UserPrefs.EnergyMax,
						EnergyFullTime:                 "2023-03-20 03:58:55",
						LicenseLiveEnergyRecoverlyTime: 60,
						EnergyFullNeedTime:             0,
						OverMaxEnergy:                  config.Conf.UserPrefs.OverMaxEnergy,
						TrainingEnergy:                 100,
						TrainingEnergyMax:              100,
						FriendMax:                      99,
						InviteCode:                     config.Conf.UserPrefs.InviteCode,
						InsertDate:                     "2015-08-10 18:58:30",
						UpdateDate:                     "2018-08-09 18:13:12",
						TutorialState:                  -1,
						DiamondCoin:                    0,
						CrystalCoin:                    0,
						LpRecoveryItem:                 []model.LpRecoveryItem{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, userInfoResp)
			default:
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "navigation":
			if v.Action == "specialCutin" {
				// key = "special_cutin_result"
				cutinResp := model.SpecialCutinResp{
					Result: model.SpecialCutinRes{
						SpecialCutinList: []any{},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, cutinResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "award":
			if v.Action == "awardInfo" {
				// key = "award_result"
				var aIdList []int
				err := ss.MainEng.Table("award_m").Cols("award_id").Find(&aIdList)
				if ss.CheckErr(err) {
					return
				}

				var aId int
				_, err = ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Cols("award_id").Get(&aId)
				if ss.CheckErr(err) {
					return
				}

				awardsList := []model.AwardInfo{}
				for _, id := range aIdList {
					isSet := false
					if id == aId {
						isSet = true
					}
					awardsList = append(awardsList, model.AwardInfo{
						AwardID:    id,
						IsSet:      isSet,
						InsertDate: time.Now().Format("2006-01-02 03:04:05"),
					})
				}

				awardResp := model.AwardInfoResp{
					Result: model.AwardInfoRes{
						AwardInfo: awardsList,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, awardResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "background":
			if v.Action == "backgroundInfo" {
				// key = "background_result"
				var bIdList []int
				err := ss.MainEng.Table("background_m").Cols("background_id").Find(&bIdList)
				if ss.CheckErr(err) {
					return
				}

				var bId int
				_, err = ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Cols("background_id").Get(&bId)
				if ss.CheckErr(err) {
					return
				}

				backgroundsList := []model.BackgroundInfo{}
				for _, id := range bIdList {
					isSet := false
					if id == bId {
						isSet = true
					}
					backgroundsList = append(backgroundsList, model.BackgroundInfo{
						BackgroundID: id,
						IsSet:        isSet,
						InsertDate:   time.Now().Format("2006-01-02 03:04:05"),
					})
				}

				backgroundResp := model.BackgroundInfoResp{
					Result: model.BackgroundInfoRes{
						BackgroundInfo: backgroundsList,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, backgroundResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "stamp":
			if v.Action == "stampInfo" {
				// key = "stamp_result"
				stampResp := honokautils.ReadAllText("assets/serverdata/stamp.json")
				var mStampResp any
				err = json.Unmarshal([]byte(stampResp), &mStampResp)
				if ss.CheckErr(err) {
					return
				}
				results = append(results, mStampResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "exchange":
			if v.Action == "owningPoint" {
				// key = "exchange_point_result"
				var exchangeIds []int
				exPointsList := []model.ExchangePointList{}
				err = ss.MainEng.Table("exchange_point_m").Cols("exchange_point_id").OrderBy("exchange_point_id ASC").Find(&exchangeIds)
				if ss.CheckErr(err) {
					return
				}

				for _, id := range exchangeIds {
					exPointsList = append(exPointsList, model.ExchangePointList{
						Rarity:        id,
						ExchangePoint: 9999,
					})
				}
				exPointsResp := model.ExchangePointResp{
					Result: model.ExchangePointRes{
						ExchangePointList: exPointsList,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, exPointsResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "livese":
			if v.Action == "liveseInfo" {
				// key = "live_se_result"
				liveSeResp := model.LiveSeInfoResp{
					Result: model.LiveSeInfoRes{
						LiveSeList: []int{1, 2, 3},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, liveSeResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "liveicon":
			if v.Action == "liveiconInfo" {
				// key = "live_icon_result"
				liveIconResp := model.LiveIconInfoResp{
					Result: model.LiveIconInfoRes{
						LiveNotesIconList: []int{1, 2, 3},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, liveIconResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "item":
			if v.Action == "list" {
				// key = "item_list_result"
				itemResp := honokautils.ReadAllText("assets/serverdata/item.json")
				var mItemResp any
				err = json.Unmarshal([]byte(itemResp), &mItemResp)
				if ss.CheckErr(err) {
					return
				}
				results = append(results, mItemResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "marathon":
			if v.Action == "marathonInfo" {
				// key = "marathon_result"
				marathonResp := model.MarathonInfoResp{
					Result:     []any{},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, marathonResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "challenge":
			if v.Action == "challengeInfo" {
				// key = "challenge_result"
				challengeResp := model.ChallengeInfoResp{
					Result:     []any{},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, challengeResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "museum":
			if v.Action == "info" {
				// key = "museum_result"
				museumRes := []MuseumRes{}
				var museumIds []int
				var smileBuff, pureBuff, coolBuff int
				err = ss.MainEng.Table("museum_contents_m").Cols("museum_contents_id,smile_buff,pure_buff,cool_buff").
					OrderBy("museum_contents_id ASC").Find(&museumRes)
				if ss.CheckErr(err) {
					return
				}

				for _, res := range museumRes {
					smileBuff += res.SmileBuff
					pureBuff += res.PureBuff
					coolBuff += res.CoolBuff
					museumIds = append(museumIds, res.MuseumContentsId)
				}
				museumInfoResp := model.MuseumInfoResp{
					Result: model.MuseumInfoRes{
						MuseumInfo: model.Museum{
							Parameter: model.MuseumParameter{
								Smile: smileBuff,
								Pure:  pureBuff,
								Cool:  coolBuff,
							},
							ContentsIDList: museumIds,
						},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, museumInfoResp)
			} else {
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		case "profile":
			switch v.Action {
			case "liveCnt":
				// key = "profile_livecnt_result"
				difficultyResp := model.DifficultyResp{
					Result: []model.DifficultyRes{
						{
							Difficulty: 1,
							ClearCnt:   315,
						},
						{
							Difficulty: 2,
							ClearCnt:   310,
						},
						{
							Difficulty: 3,
							ClearCnt:   314,
						},
						{
							Difficulty: 4,
							ClearCnt:   455,
						},
						{
							Difficulty: 6,
							ClearCnt:   233,
						},
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, difficultyResp)
			case "cardRanking":
				// key = "profile_card_ranking_result"
				var result []any
				love := honokautils.ReadAllText("assets/serverdata/love.json")
				err := json.Unmarshal([]byte(love), &result)
				if ss.CheckErr(err) {
					return
				}

				loveResp := model.LoveResp{
					Result:     result,
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, loveResp)
			case "profileInfo":
				// key = "profile_info_result"
				pref := tools.UserPref{}
				exists, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Get(&pref)
				if ss.CheckErr(err) {
					return
				}
				if !exists {
					ctx.String(http.StatusForbidden, ErrorMsg)
					return
				}

				commonUnit, err := ss.MainEng.Table("common_unit_m").Count()
				if ss.CheckErr(err) {
					return
				}
				userUnit, err := ss.UserEng.Table("user_unit_m").Where("user_id = ?", ctx.GetString("userid")).Count()
				if ss.CheckErr(err) {
					return
				}

				unitData := model.UnitData{}
				exists, err = ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", pref.UnitOwningUserID).Get(&unitData)
				if ss.CheckErr(err) {
					return
				}

				isCommon := true
				if !exists {
					_, err = ss.UserEng.Table("user_unit_m").
						Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ctx.GetString("userid")).Get(&unitData)
					if ss.CheckErr(err) {
						return
					}
					isCommon = false
				}

				var attrId, maxHp, baseSmile, basePure, baseCool int
				var smileMax, pureMax, coolMax int
				if isCommon {
					// 公共卡片仅为100级属性
					_, err = ss.MainEng.Table("unit_m").Where("unit_id = ?", unitData.UnitID).
						Cols("attribute_id,hp_max,smile_max,pure_max,cool_max").Get(&attrId, &maxHp, &baseSmile, &basePure, &baseCool)
					utils.CheckErr(err)

					// 偷懒起见不计算饰品、宝石、回忆画廊等属性加成
					smileMax = baseSmile
					pureMax = basePure
					coolMax = baseCool
					// } else {
					// 	// 用户卡片需要根据等级计算属性
					// 	// TODO
				}

				var accessoryOwningId, accessoryId, exp int
				_, err = ss.UserEng.Table("accessory_wear_m").Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ctx.GetString("userid")).
					Cols("accessory_owning_user_id").Get(&accessoryOwningId)
				if ss.CheckErr(err) {
					return
				}
				_, err = ss.MainEng.Table("common_accessory_m").Where("accessory_owning_user_id = ?", accessoryOwningId).
					Cols("accessory_id,exp").Get(&accessoryId, &exp)
				if ss.CheckErr(err) {
					return
				}
				accessoryInfo := model.AccessoryInfo{
					AccessoryOwningUserID: accessoryOwningId,
					AccessoryID:           accessoryId,
					Exp:                   exp,
					NextExp:               0,
					Level:                 8,
					MaxLevel:              8,
					RankUpCount:           4,
					FavoriteFlag:          true,
				}

				removeSkillIds := []int{}
				err = ss.UserEng.Table("skill_equip_m").Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ctx.GetString("userid")).
					Cols("unit_removable_skill_id").Find(&removeSkillIds)
				if ss.CheckErr(err) {
					return
				}

				userId, err := strconv.Atoi(config.Conf.UserPrefs.InviteCode)
				if ss.CheckErr(err) {
					return
				}

				profileResp := model.ProfileResp{
					Result: model.ProfileRes{
						UserInfo: model.ProfileUserInfo{
							UserID:               userId,
							Name:                 pref.UserName,
							Level:                config.Conf.UserPrefs.Level,
							CostMax:              100,
							UnitMax:              5000,
							EnergyMax:            config.Conf.UserPrefs.EnergyMax,
							FriendMax:            99,
							UnitCnt:              int(commonUnit + userUnit),
							InviteCode:           config.Conf.UserPrefs.InviteCode,
							ElapsedTimeFromLogin: "14\u5c0f\u65f6\u524d",
							Introduction:         pref.UserDesc,
						},
						CenterUnitInfo: model.CenterUnitInfo{
							UnitOwningUserID:           unitData.UnitOwningUserID,
							UnitID:                     unitData.UnitID,
							Exp:                        unitData.Exp,
							NextExp:                    unitData.NextExp,
							Level:                      unitData.Level,
							LevelLimitID:               unitData.LevelLimitID,
							MaxLevel:                   unitData.MaxLevel,
							Rank:                       unitData.Rank,
							MaxRank:                    unitData.MaxRank,
							Love:                       unitData.Love,
							MaxLove:                    unitData.MaxLove,
							UnitSkillLevel:             unitData.UnitSkillLevel,
							MaxHp:                      unitData.MaxHp,
							FavoriteFlag:               unitData.FavoriteFlag,
							DisplayRank:                unitData.DisplayRank,
							UnitSkillExp:               unitData.UnitSkillExp,
							UnitRemovableSkillCapacity: unitData.UnitRemovableSkillCapacity,
							Attribute:                  attrId,
							Smile:                      baseSmile,
							Cute:                       basePure,
							Cool:                       baseCool,
							IsLoveMax:                  unitData.IsLoveMax,
							IsLevelMax:                 unitData.IsLevelMax,
							IsRankMax:                  unitData.IsRankMax,
							IsSigned:                   unitData.IsSigned,
							IsSkillLevelMax:            unitData.IsSkillLevelMax,
							SettingAwardID:             pref.AwardID,
							RemovableSkillIds:          removeSkillIds,
							AccessoryInfo:              accessoryInfo,
							Costume:                    model.Costume{},
							TotalSmile:                 smileMax,
							TotalCute:                  pureMax,
							TotalCool:                  coolMax,
							TotalHp:                    maxHp,
						},
						NaviUnitInfo: model.NaviUnitInfo{
							UnitOwningUserID:            unitData.UnitOwningUserID,
							UnitID:                      unitData.UnitID,
							Exp:                         unitData.Exp,
							NextExp:                     unitData.NextExp,
							Level:                       unitData.Level,
							MaxLevel:                    unitData.MaxLevel,
							LevelLimitID:                unitData.LevelLimitID,
							Rank:                        unitData.Rank,
							MaxRank:                     unitData.MaxRank,
							Love:                        unitData.Love,
							MaxLove:                     unitData.MaxLove,
							UnitSkillExp:                unitData.UnitSkillExp,
							UnitSkillLevel:              unitData.UnitSkillLevel,
							MaxHp:                       unitData.MaxHp,
							UnitRemovableSkillCapacity:  unitData.UnitRemovableSkillCapacity,
							FavoriteFlag:                unitData.FavoriteFlag,
							DisplayRank:                 unitData.DisplayRank,
							IsLoveMax:                   unitData.IsLoveMax,
							IsLevelMax:                  unitData.IsLevelMax,
							IsRankMax:                   unitData.IsRankMax,
							IsSigned:                    unitData.IsSigned,
							IsSkillLevelMax:             unitData.IsSkillLevelMax,
							IsRemovableSkillCapacityMax: unitData.IsRemovableSkillCapacityMax,
							InsertDate:                  "2016-10-11 10:33:03",
							TotalSmile:                  smileMax,
							TotalCute:                   pureMax,
							TotalCool:                   coolMax,
							TotalHp:                     maxHp,
							RemovableSkillIds:           removeSkillIds,
						},
						IsAlliance:          false,
						FriendStatus:        0,
						SettingAwardID:      113,
						SettingBackgroundID: 143,
					},
					Status:     200,
					CommandNum: false,
					TimeStamp:  time.Now().Unix(),
				}
				results = append(results, profileResp)
			default:
				err = fmt.Errorf("unimplemented: %s", v.Module+":"+v.Action)
			}
		default:
			err = fmt.Errorf("unimplemented: %s", v.Module)
		}
	}

	if ss.CheckErr(err) {
		return
	}

	apiResp := model.ApiResp{
		ResponseData: results,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(apiResp)
}

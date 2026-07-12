package secretbox

import (
	"fmt"
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	"honoka-chan/internal/router"
	secretboxapischema "honoka-chan/internal/schema/api/secretbox"
	secretboxschema "honoka-chan/internal/schema/secretbox"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func pon(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	ponReq := secretboxschema.PonReq{}
	err := honokautils.ParseRequestData(ctx, &ponReq)
	if ss.CheckErr(err) {
		return
	}

	data, err := drawData(ctx, ponReq.SecretBoxID, ponReq.ID)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(secretboxschema.PonResp{
		ResponseData: data,
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/secretbox/pon", middleware.Common, pon)
}

func drawData(ctx *gin.Context, secretBoxID, costID int) (secretboxschema.PonData, error) {
	if secretBoxID <= 0 || costID <= 0 {
		return secretboxschema.PonData{}, fmt.Errorf("invalid secretbox draw params: secret_box_id=%d id=%d", secretBoxID, costID)
	}

	ss := session.Get(ctx)
	config, err := loadSecretBoxAllData()
	if err != nil {
		return secretboxschema.PonData{}, err
	}

	memberCategory, page, err := findSecretBoxPage(config, secretBoxID)
	if err != nil {
		return secretboxschema.PonData{}, err
	}

	button, cost, err := findSecretBoxCost(page, costID)
	if err != nil {
		return secretboxschema.PonData{}, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	drawnUnits, lowestRarity, err := buildDrawUnits(ss, rng, memberCategory, page, button, cost)
	if err != nil {
		return secretboxschema.PonData{}, err
	}

	for i := range page.ButtonList {
		for j := range page.ButtonList[i].CostList {
			page.ButtonList[i].CostList[j].Payable = true
		}
	}
	page.SecretBoxInfo.PonCount += cost.UnitCount

	itemList := config.Result.ItemList
	for i := range itemList {
		if itemList[i].Amount < 999 {
			itemList[i].Amount = 999
		}
	}

	return secretboxschema.PonData{
		IsUnitMax: false,
		ItemList:  itemList,
		GaugeInfo: secretboxapischema.GaugeInfo{
			MaxGaugePoint: 100,
			GaugePoint:    0,
		},
		ButtonList:               page.ButtonList,
		SecretBoxInfo:            page.SecretBoxInfo,
		SecretBoxItems:           secretboxschema.SecretBoxItems{Unit: drawnUnits, Item: []secretboxschema.SecretBoxRewardItem{}},
		BeforeUserInfo:           ss.GetUserInfo(),
		AfterUserInfo:            ss.GetUserInfo(),
		FreeMuseGachaFlag:        false,
		FreeAqoursGachaFlag:      false,
		LowestRarity:             lowestRarity,
		PromotionPerformanceRate: 25,
		SecretBoxParcelType:      2,
		LimitBonusInfo:           []any{},
		LimitBonusRewards:        []any{},
		UnitSupportList:          []any{},
	}, nil
}

func buildDrawUnits(ss *session.Session, rng *rand.Rand, memberCategory int, page secretboxapischema.SecretBoxPage, button secretboxapischema.SecretBoxButton, cost secretboxapischema.SecretBoxCost) ([]secretboxschema.SecretBoxUnitItem, int, error) {
	drawCount := cost.UnitCount
	if drawCount <= 0 {
		drawCount = 1
	}

	rarities := make([]int, 0, drawCount)
	for i := 0; i < drawCount; i++ {
		rarities = append(rarities, pickRarity(rng, page, button, i == drawCount-1 && drawCount > 1))
	}

	if drawCount > 1 {
		guaranteedRarity := raritySR
		if hasGuaranteedUR(button, cost) {
			guaranteedRarity = rarityUR
		}

		guaranteed := false
		for _, rarity := range rarities {
			if satisfiesGuaranteedRarity(rarity, guaranteedRarity) {
				guaranteed = true
				break
			}
		}
		if !guaranteed {
			if guaranteedRarity == rarityUR {
				rarities[drawCount-1] = rarityUR
			} else {
				rarities[drawCount-1] = pickGuaranteedRarity(rng, page, button)
			}
		}
	}

	result := make([]secretboxschema.SecretBoxUnitItem, 0, drawCount)
	lowestRarity := 0
	for _, rarity := range rarities {
		unitID, err := randomUnitID(ss, rng, memberCategory, page, button, rarity)
		if err != nil {
			return nil, 0, err
		}

		unitData, err := getUserUnitByUnitID(ss, unitID)
		if err != nil {
			return nil, 0, err
		}

		item := buildUnitItem(unitData)
		item.IsHit = shouldShowDirectURHit(rng, page, button, rarity)
		result = append(result, item)
		if lowestRarity == 0 || isLowerQualityRarity(rarity, lowestRarity) {
			lowestRarity = rarity
		}
	}

	return result, lowestRarity, nil
}

func buildUnitItem(unitData unitmodel.UnitDataMap) secretboxschema.SecretBoxUnitItem {
	isSupportMember := unitData.MaxLove == 0 && unitData.UnitSkillLevel == 0 && unitData.UnitRemovableSkillCapacity == 0
	insertDate := time.Unix(unitData.InsertDate, 0).Format("2006-01-02 15:04:05")

	return secretboxschema.SecretBoxUnitItem{
		UnitOwningUserID:            unitData.UnitOwningUserID,
		UnitOwningIDs:               []int{unitData.UnitOwningUserID},
		UnitID:                      unitData.UnitID,
		Exp:                         unitData.Exp,
		NextExp:                     0,
		Level:                       unitData.Level,
		MaxLevel:                    unitData.MaxLevel,
		LevelLimitID:                unitData.LevelLimitID,
		Rank:                        unitData.Rank,
		MaxRank:                     unitData.MaxRank,
		Love:                        unitData.Love,
		MaxLove:                     unitData.MaxLove,
		UnitSkillLevel:              unitData.UnitSkillLevel,
		SkillLevel:                  unitData.UnitSkillLevel,
		UnitSkillExp:                unitData.UnitSkillExp,
		MaxHp:                       unitData.MaxHp,
		FavoriteFlag:                unitData.FavoriteFlag,
		DisplayRank:                 unitData.DisplayRank,
		UnitRemovableSkillCapacity:  unitData.UnitRemovableSkillCapacity,
		MaxRemovableSkillCapacity:   unitData.UnitRemovableSkillCapacity,
		Attribute:                   unitData.Attribute,
		Smile:                       unitData.Smile,
		Cute:                        unitData.Cute,
		Cool:                        unitData.Cool,
		IsRankMax:                   unitData.IsRankMax,
		IsLoveMax:                   unitData.IsLoveMax,
		IsLevelMax:                  unitData.IsLevelMax,
		IsSigned:                    unitData.IsSigned,
		IsSkillLevelMax:             unitData.IsSkillLevelMax,
		IsRemovableSkillCapacityMax: unitData.IsRemovableSkillCapacityMax,
		IsSupportMember:             isSupportMember,
		NewUnitFlag:                 false,
		UnitRarityID:                unitData.Rarity,
		IsHit:                       unitData.Rarity == rarityUR,
		RemovableSkillIDs:           []int{},
		InsertDate:                  insertDate,
		TotalSmile:                  unitData.Smile,
		TotalCute:                   unitData.Cute,
		TotalCool:                   unitData.Cool,
		TotalHp:                     unitData.MaxHp,
	}
}

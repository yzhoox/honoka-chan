package api

import (
	"honoka-chan/internal/constant"
	"honoka-chan/internal/handler/api/achievement"
	"honoka-chan/internal/handler/api/album"
	"honoka-chan/internal/handler/api/award"
	"honoka-chan/internal/handler/api/background"
	"honoka-chan/internal/handler/api/banner"
	"honoka-chan/internal/handler/api/challenge"
	"honoka-chan/internal/handler/api/costume"
	"honoka-chan/internal/handler/api/eventscenario"
	"honoka-chan/internal/handler/api/exchange"
	"honoka-chan/internal/handler/api/item"
	"honoka-chan/internal/handler/api/live"
	"honoka-chan/internal/handler/api/liveicon"
	"honoka-chan/internal/handler/api/livese"
	"honoka-chan/internal/handler/api/login"
	"honoka-chan/internal/handler/api/marathon"
	"honoka-chan/internal/handler/api/multiunit"
	"honoka-chan/internal/handler/api/museum"
	"honoka-chan/internal/handler/api/navigation"
	"honoka-chan/internal/handler/api/notice"
	"honoka-chan/internal/handler/api/payment"
	"honoka-chan/internal/handler/api/profile"
	"honoka-chan/internal/handler/api/reward"
	"honoka-chan/internal/handler/api/scenario"
	"honoka-chan/internal/handler/api/secretbox"
	"honoka-chan/internal/handler/api/stamp"
	"honoka-chan/internal/handler/api/subscenario"
	"honoka-chan/internal/handler/api/unit"
	"honoka-chan/internal/handler/api/user"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	apischema "honoka-chan/internal/schema/api"
	commonschema "honoka-chan/internal/schema/common"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func api(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	apiReq := []apischema.ApiReq{}
	err := honokautils.ParseRequestData(ctx, &apiReq)
	if ss.CheckErr(err) {
		return
	}

	var result any
	var results []any
	for _, v := range apiReq {
		// fmt.Println(v.Module, v.Action)
		switch v.Module {
		case "album":
			result, err = album.AlbumApi(ctx, v.Action)
		case "achievement":
			result, err = achievement.AchievementApi(ctx, v.Action)
		case "award":
			result, err = award.AwardApi(ctx, v.Action)
		case "background":
			result, err = background.BackgroundApi(ctx, v.Action)
		case "banner":
			result, err = banner.BannerApi(v.Action)
		case "challenge":
			result, err = challenge.ChallengeApi(v.Action)
		case "costume":
			result, err = costume.CostumeApi(v.Action)
		case "eventscenario":
			result, err = eventscenario.EventScenarioApi(ctx, v.Action)
		case "exchange":
			result, err = exchange.ExchangeApi(ctx, v.Action)
		case "item":
			result, err = item.ItemApi(v.Action)
		case "live":
			result, err = live.LiveApi(ctx, v.Action)
		case "liveicon":
			result, err = liveicon.LiveIconApi(v.Action)
		case "livese":
			result, err = livese.LiveSeApi(v.Action)
		case "login":
			result, err = login.LoginApi(ctx, v.Action)
		case "marathon":
			result, err = marathon.MarathonApi(v.Action)
		case "multiunit":
			result, err = multiunit.MultiUnitApi(ctx, v.Action)
		case "museum":
			result, err = museum.MuseumApi(ctx, v.Action)
		case "navigation":
			result, err = navigation.NavigationApi(v.Action)
		case "notice":
			result, err = notice.NoticeApi(v.Action)
		case "payment":
			result, err = payment.PaymentApi(ctx, v.Action)
		case "profile":
			result, err = profile.ProfileApi(ctx, v.Action, v.UserID)
		case "reward":
			result, err = reward.RewardApi(v.Action)
		case "scenario":
			result, err = scenario.ScenarioApi(ctx, v.Action)
		case "secretbox":
			result, err = secretbox.SecretBoxApi(v.Action)
		case "stamp":
			result, err = stamp.StampApi(v.Action)
		case "subscenario":
			result, err = subscenario.SubscenarioApi(ctx, v.Action)
		case "unit":
			result, err = unit.UnitApi(ctx, v.Action)
		case "user":
			result, err = user.UserApi(ctx, v.Action)
		default:
			err = honokautils.NewUnimplementedModuleError(v.Module)
		}

		if honokautils.IsUnimplementedError(err) {
			results = append(results, commonschema.ApiErrorResp{
				Result: commonschema.ErrorData{
					ErrorCode: constant.ErrorCodeUnknown,
				},
				Status:     600,
				CommandNum: false,
				TimeStamp:  time.Now().Unix(),
			})
			continue
		}
		if ss.CheckErr(err) {
			return
		}
		results = append(results, result)
	}

	apiResp := apischema.ApiResp{
		ResponseData: results,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(apiResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/api", middleware.Common, api)
}

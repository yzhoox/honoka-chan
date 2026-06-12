package achievement

import (
	apiachievement "honoka-chan/internal/handler/api/achievement"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	achievementschema "honoka-chan/internal/schema/achievement"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

const pagingAccomplishedPageSize = 10

func pagingAccomplishedList(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := achievementschema.PagingAccomplishedListReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	grouped, err := apiachievement.ListVisibleAccomplishedAchievements(ss)
	if ss.CheckErr(err) {
		return
	}

	items := grouped[req.FilterCategoryID]
	start := min(max(req.FromCount, 0), len(items))

	end := min(start+pagingAccomplishedPageSize, len(items))

	ss.Respond(achievementschema.PagingAccomplishedListResp{
		ResponseData: items[start:end],
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/achievement/pagingAccomplishedList", middleware.Common, pagingAccomplishedList)
}

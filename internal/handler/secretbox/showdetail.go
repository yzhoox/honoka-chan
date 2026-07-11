package secretbox

import (
	"fmt"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	secretboxschema "honoka-chan/internal/schema/secretbox"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func showDetail(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	showDetailReq := secretboxschema.ShowDetailReq{}
	err := honokautils.ParseRequestData(ctx, &showDetailReq)
	if ss.CheckErr(err) {
		return
	}

	data, err := showDetailData(showDetailReq.SecretBoxID)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(secretboxschema.ShowDetailResp{
		ResponseData: data,
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/secretbox/showDetail", middleware.Common, showDetail)
}

func showDetailData(secretBoxID int) (secretboxschema.ShowDetailData, error) {
	if secretBoxID <= 0 {
		return secretboxschema.ShowDetailData{}, fmt.Errorf("invalid secret_box_id: %d", secretBoxID)
	}

	config, err := loadSecretBoxAllData()
	if err != nil {
		return secretboxschema.ShowDetailData{}, err
	}

	memberCategory, page, err := findSecretBoxPage(config, secretBoxID)
	if err != nil {
		return secretboxschema.ShowDetailData{}, err
	}

	result := secretboxschema.ShowDetailData{
		URL: fmt.Sprintf("/webview.php/secretbox/detail?id=%d&no_title=1", secretBoxID),
	}

	switch page.SecretBoxInfo.SecretBoxType {
	case 5:
		buttons := make([]secretboxschema.ButtonTypeUnitLineUp, 0, len(page.ButtonList))
		for _, button := range page.ButtonList {
			buttons = append(buttons, secretboxschema.ButtonTypeUnitLineUp{
				SecretBoxButtonType: button.SecretBoxButtonType,
				SecretBoxName:       button.SecretBoxName,
				UnitLineUp:          buildUnitLineUp(memberCategory, page, button),
			})
		}
		result.ButtonTypeUnitLineUp = buttons
	default:
		if len(page.ButtonList) == 0 {
			return secretboxschema.ShowDetailData{}, fmt.Errorf("secretbox %d has no buttons", secretBoxID)
		}
		result.UnitLineUp = buildUnitLineUp(memberCategory, page, page.ButtonList[0])
	}

	return result, nil
}

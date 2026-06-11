package system

import (
	"encoding/json"
	"errors"
	"honoka-chan/config"
	"honoka-chan/internal/router"
	systemschema "honoka-chan/internal/schema/system"
	"honoka-chan/pkg/db"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	startedAt    = time.Now()
	lastReloadAt atomic.Int64
)

const dateTimeFormat = "2006-01-02 15:04:05"

func MarkStarted(t time.Time) {
	startedAt = t
}

func Health() (systemschema.HealthResp, int) {
	appName := config.Conf.AppName
	if appName == "" {
		appName = "honoka-chan"
	}

	mainDBStatus := "not initialized"
	if db.MainEng != nil {
		mainDBStatus = "ok"
		if err := db.MainEng.Ping(); err != nil {
			mainDBStatus = err.Error()
		}
	}

	userDBStatus := "not initialized"
	if db.UserEng != nil {
		userDBStatus = "ok"
		if err := db.UserEng.Ping(); err != nil {
			userDBStatus = err.Error()
		}
	}

	status := "ok"
	statusCode := http.StatusOK
	if db.MainEng == nil || db.UserEng == nil {
		status = "stopped"
	}

	resp := systemschema.HealthResp{
		Status:                status,
		AppName:               appName,
		Version:               config.PackageVersion,
		StartedAt:             startedAt.Format(dateTimeFormat),
		UptimeSeconds:         int64(time.Since(startedAt).Seconds()),
		ListenPort:            config.Conf.Settings.ListenPort,
		CdnServer:             config.Conf.Settings.CdnServer,
		ReloadTokenConfigured: strings.TrimSpace(config.Conf.Settings.ReloadToken) != "",
		MainDB:                mainDBStatus,
		UserDB:                userDBStatus,
	}

	if ts := lastReloadAt.Load(); ts > 0 {
		resp.LastReloadAt = time.Unix(ts, 0).Format(dateTimeFormat)
	}

	if status == "ok" && (mainDBStatus != "ok" || userDBStatus != "ok") {
		resp.Status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	return resp, statusCode
}

func HealthJSON() string {
	resp, _ := Health()
	data, err := json.Marshal(resp)
	if err != nil {
		return `{"status":"error","app_name":"honoka-chan"}`
	}
	return string(data)
}

func Reload(token string, requireToken bool) (systemschema.ReloadResp, int, error) {
	expectedToken := strings.TrimSpace(config.Conf.Settings.ReloadToken)
	if requireToken && expectedToken == "" {
		return systemschema.ReloadResp{
			Status:  "error",
			Message: "reload_token is not configured",
		}, http.StatusServiceUnavailable, errors.New("reload_token is not configured")
	}

	if requireToken && strings.TrimSpace(token) != expectedToken {
		return systemschema.ReloadResp{
			Status:  "error",
			Message: "invalid reload token",
		}, http.StatusUnauthorized, errors.New("invalid reload token")
	}

	config.ReloadConfig()
	now := time.Now()
	lastReloadAt.Store(now.Unix())

	return systemschema.ReloadResp{
		Status:     "ok",
		Message:    "configuration reloaded",
		ReloadedAt: now.Format(dateTimeFormat),
	}, http.StatusOK, nil
}

func health(ctx *gin.Context) {
	resp, statusCode := Health()
	ctx.JSON(statusCode, resp)
}

func reload(ctx *gin.Context) {
	token := strings.TrimSpace(ctx.GetHeader("X-Reload-Token"))
	if token == "" {
		token = strings.TrimSpace(ctx.Query("token"))
	}
	if token == "" {
		token = strings.TrimSpace(ctx.PostForm("token"))
	}

	resp, statusCode, _ := Reload(token, true)
	ctx.JSON(statusCode, resp)
}

func init() {
	router.AddHandler("/", "GET", "/system/health", health)
	router.AddHandler("/", "POST", "/system/reload", reload)
}

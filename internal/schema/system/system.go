package systemschema

type HealthResp struct {
	Status                string `json:"status"`
	AppName               string `json:"app_name"`
	Version               string `json:"version"`
	StartedAt             string `json:"started_at"`
	UptimeSeconds         int64  `json:"uptime_seconds"`
	LastReloadAt          string `json:"last_reload_at,omitempty"`
	ListenPort            string `json:"listen_port"`
	CdnServer             string `json:"cdn_server"`
	ReloadTokenConfigured bool   `json:"reload_token_configured"`
	MainDB                string `json:"main_db"`
	UserDB                string `json:"user_db"`
}

type ReloadResp struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	ReloadedAt string `json:"reloaded_at,omitempty"`
}

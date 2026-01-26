package ghomeschema

type ActiveData struct {
	Message string `json:"message"`
	Result  int    `json:"result"`
}

type ActiveResp struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data ActiveData `json:"data"`
}

type InitializeData struct {
	BrandLogo                 string   `json:"brand_logo"`
	BrandName                 string   `json:"brand_name"`
	DaoyuClientid             string   `json:"daoyu_clientid"`
	DaoyuDownloadURL          string   `json:"daoyu_download_url"`
	DeviceFeature             string   `json:"device_feature"`
	DisplayThirdaccout        int      `json:"display_thirdaccout"`
	ForceShowAgreement        int      `json:"force_show_agreement"`
	GreportLogLevel           string   `json:"greport_log_level"`
	GuestEnable               int      `json:"guest_enable"`
	IsMatch                   int      `json:"is_match"`
	LogLevel                  string   `json:"log_level"`
	LoginButton               []string `json:"login_button"`
	LoginIcon                 []any    `json:"login_icon"`
	LoginLimitEnable          int      `json:"login_limit_enable"`
	NeedFloatWindowPermission int      `json:"need_float_window_permission"`
	NewDeviceIDServer         string   `json:"new_device_id_server"`
	QqAppID                   string   `json:"qq_appId"`
	QqKey                     string   `json:"qq_key"`
	ShowGuestConfirm          int      `json:"show_guest_confirm"`
	VoicetipButton            int      `json:"voicetip_button"`
	VoicetipOne               string   `json:"voicetip_one"`
	VoicetipTwo               string   `json:"voicetip_two"`
	WegameAppid               string   `json:"wegame_appid"`
	WegameAppkey              string   `json:"wegame_appkey"`
	WegameClientid            string   `json:"wegame_clientid"`
	WegameCompanyID           string   `json:"wegame_companyId"`
	WegameLoginURL            string   `json:"wegame_loginUrl"`
	WeiboAppKey               string   `json:"weibo_appKey"`
	WeiboRedirectURL          string   `json:"weibo_redirectUrl"`
	WeixinAppID               string   `json:"weixin_appId"`
	WeixinKey                 string   `json:"weixin_key"`
}

type InitializeResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type LoginData struct {
	Activation           int    `json:"activation"`
	Autokey              string `json:"autokey"`
	CaptchaParams        string `json:"captchaParams"`
	CheckCodeGUID        string `json:"checkCodeGuid"`
	CheckCodeURL         string `json:"checkCodeUrl"`
	HasExtendAccs        int    `json:"hasExtendAccs"`
	HasRealInfo          int    `json:"has_realInfo"`
	ImagecodeType        int    `json:"imagecodeType"`
	IsNewUser            int    `json:"isNewUser"`
	Message              string `json:"message"`
	NextAction           int    `json:"nextAction"`
	PromptMsg            string `json:"prompt_msg"`
	RealInfoNotification string `json:"realInfoNotification"`
	RealInfoForce        int    `json:"realInfo_force"`
	RealInfoForcePay     int    `json:"realInfo_force_pay"`
	RealInfoStatus       int    `json:"realInfo_status"`
	RealInfoStatusPay    int    `json:"realInfo_status_pay"`
	Result               int    `json:"result"`
	SdgHeight            int    `json:"sdg_height"`
	SdgWidth             int    `json:"sdg_width"`
	Ticket               string `json:"ticket"`
	UserAttribute        string `json:"userAttribute"`
	UserID               int    `json:"userid"`
}

type LoginResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type LoginAutoData struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
	Autokey string `json:"autokey,omitempty"`
	UserId  string `json:"userid,omitempty"`
	Ticket  string `json:"ticket,omitempty"`
}

type LoginAutoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type ReportRoleResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

package ghomeschema

type AgreementData struct {
}

type AgreementResp struct {
	ReturnCode    int           `json:"return_code"`
	ErrorType     int           `json:"error_type"`
	ReturnMessage string        `json:"return_message"`
	Data          AgreementData `json:"data"`
}

type AppReportData struct {
	NeedReport int `json:"needReport"`
}

type AppReportResp struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data AppReportData `json:"data"`
}

package liveicon

type InfoData struct {
	LiveNotesIconList []int `json:"live_notes_icon_list"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}

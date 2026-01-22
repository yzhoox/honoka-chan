package museum

type Parameter struct {
	Smile int `json:"smile"`
	Pure  int `json:"pure"`
	Cool  int `json:"cool"`
}

type Info struct {
	Parameter      Parameter `json:"parameter"`
	ContentsIDList []int     `json:"contents_id_list"`
}

type InfoData struct {
	MuseumInfo Info `json:"museum_info"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}

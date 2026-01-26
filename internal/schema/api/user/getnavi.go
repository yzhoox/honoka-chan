package userapischema

type User struct {
	UserID           int `json:"user_id"`
	UnitOwningUserID int `json:"unit_owning_user_id"`
}

type GetNaviData struct {
	User User `json:"user"`
}

type GetNaviResp struct {
	Result     GetNaviData `json:"result"`
	Status     int         `json:"status"`
	CommandNum bool        `json:"commandNum"`
	TimeStamp  int64       `json:"timeStamp"`
}

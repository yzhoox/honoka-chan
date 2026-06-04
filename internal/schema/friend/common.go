package friendschema

type UserIDReq struct {
	UserID int `json:"user_id"`
}

type ResponseReq struct {
	UserID int `json:"user_id"`
	Status int `json:"status"`
}

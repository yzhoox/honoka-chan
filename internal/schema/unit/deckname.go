package unitschema

type DeckNameReq struct {
	Module     string `json:"module"`
	UnitDeckID int    `json:"unit_deck_id"`
	Action     string `json:"action"`
	TimeStamp  int    `json:"timeStamp"`
	Mgd        int    `json:"mgd"`
	CommandNum string `json:"commandNum"`
	DeckName   string `json:"deck_name"`
}

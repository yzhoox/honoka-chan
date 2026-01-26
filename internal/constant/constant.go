package constant

type LiveGoalType int

const (
	LiveGoalTypeScore LiveGoalType = iota + 1
	LiveGoalTypeCombo
	LiveGoalTypeClear
)

type PreciseScoreUpdateType int

const (
	PreciseScoreUpdateTypePerfect PreciseScoreUpdateType = iota + 1
	PreciseScoreUpdateTypeScore
	PreciseScoreUpdateTypeAlways
)

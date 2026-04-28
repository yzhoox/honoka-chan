package usermodel

import "honoka-chan/internal/constant"

type UserLiveGoal struct {
	ID               int                   `xorm:"id pk autoincr"`
	LiveDifficultyID int                   `xorm:"live_difficulty_id"`
	UserID           int                   `xorm:"user_id"`
	GoalType         constant.LiveGoalType `xorm:"goal_type"`
	CompletedAt      int64                 `xorm:"completed_at"`
}

func (UserLiveGoal) TableName() string {
	return "user_live_goal"
}

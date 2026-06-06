package usermodel

import "honoka-chan/internal/constant"

type UserLiveGoal struct {
	ID               int                   `xorm:"id pk autoincr"`
	LiveDifficultyID int                   `xorm:"live_difficulty_id index"`
	UserID           int                   `xorm:"user_id index"`
	LiveGoalRewardID int                   `xorm:"live_goal_reward_id index"`
	GoalType         constant.LiveGoalType `xorm:"goal_type"`
	Rank             int                   `xorm:"rank"`
	CompletedAt      int64                 `xorm:"completed_at"`
}

func (UserLiveGoal) TableName() string {
	return "user_live_goal"
}

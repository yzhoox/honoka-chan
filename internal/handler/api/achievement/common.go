package achievement

import (
	"fmt"
	achievementapischema "honoka-chan/internal/schema/api/achievement"
	"honoka-chan/internal/session"
	"sort"
	"strings"
	"time"
)

type achievementRow struct {
	AchievementID               int `xorm:"achievement_id"`
	AchievementFilterCategoryID int `xorm:"achievement_filter_category_id"`
	DisplayFlag                 int `xorm:"display_flag"`
	StartDate                   string
	EndDate                     *string
}

func ListVisibleAccomplishedAchievements(ss *session.Session) (map[int][]achievementapischema.AchievementListItem, error) {
	var achievementList []achievementRow
	err := ss.MainEng.Table("achievement_m").
		Cols("achievement_id,achievement_filter_category_id,display_flag,start_date,end_date").
		Where("display_flag = ?", 1).
		Find(&achievementList)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	grouped := make(map[int][]achievementapischema.AchievementListItem, len(achievementFilterCategoryList))
	for _, achievement := range achievementList {
		if !isAchievementVisible(now, achievement.StartDate, achievement.EndDate) {
			continue
		}

		grouped[achievement.AchievementFilterCategoryID] = append(
			grouped[achievement.AchievementFilterCategoryID],
			achievementapischema.AchievementListItem{
				AchievementID:  achievement.AchievementID,
				Count:          1,
				IsAccomplished: true,
				InsertDate:     normalizeAchievementDate(achievement.StartDate),
				EndDate:        normalizeAchievementDatePtr(achievement.EndDate),
				RemainingTime:  "",
				IsNew:          false,
				ForDisplay:     achievement.DisplayFlag != 0,
				RewardList: []achievementapischema.AchievementRewardItem{
					{
						AddType: 3001,
						ItemID:  4,
						Amount:  1,
					},
				},
			},
		)
	}

	for filterCategoryID, items := range grouped {
		sort.Slice(items, func(i, j int) bool {
			return items[i].InsertDate > items[j].InsertDate
		})
		grouped[filterCategoryID] = items
	}

	return grouped, nil
}

func normalizeAchievementDate(value string) string {
	if value == "" {
		return "1970-01-01 00:00:00"
	}

	return strings.ReplaceAll(value, "/", "-")
}

func normalizeAchievementDatePtr(value *string) *string {
	if value == nil || *value == "" {
		return nil
	}

	normalized := normalizeAchievementDate(*value)
	return &normalized
}

func isAchievementVisible(now time.Time, startDate string, endDate *string) bool {
	start, err := parseAchievementTime(startDate)
	if err == nil && now.Before(start) {
		return false
	}

	if endDate == nil || *endDate == "" {
		return true
	}

	end, err := parseAchievementTime(*endDate)
	if err != nil {
		return true
	}

	return !now.After(end)
}

func parseAchievementTime(value string) (time.Time, error) {
	layouts := []string{
		"2006/01/02 15:04:05",
		"2006/1/2 15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid achievement time: %s", value)
}

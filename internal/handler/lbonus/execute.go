package lbonus

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/lbonus"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func execute(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	weeks := map[string]int{
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	// 本月日历
	y, m, d := time.Now().Local().Date()
	cm := m

	d1 := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	// fmt.Println(d1)
	// fmt.Println(weeks[d1.Weekday().String()])

	d2 := d1.AddDate(0, 1, -1)
	// fmt.Println(d2)

	weeksList := []lbonus.Day{}
	for c := d1; ; c = c.AddDate(0, 0, 1) {
		_, _, rd := c.Date()
		received := false
		if rd <= d {
			received = true
		}
		rw := weeks[c.Weekday().String()]
		weeksList = append(weeksList, lbonus.Day{
			Day:               rd,
			DayOfTheWeek:      rw,
			SpecialDay:        false,
			SpecialImageAsset: "",
			Received:          received,
			AdReceived:        false,
			Item: lbonus.Item{
				ItemID:  4,
				AddType: 3001,
				Amount:  1,
			},
		})
		if c.Equal(d2) {
			break
		}
	}

	// 下月日历
	y, m, _ = time.Now().AddDate(0, 1, 0).Date()
	// fmt.Println(y, m, d)

	d1 = time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	// fmt.Println(d1)
	// fmt.Println(weeks[d1.Weekday().String()])

	d2 = d1.AddDate(0, 1, -1)
	// fmt.Println(d2)

	nextWeeksList := []lbonus.Day{}
	for c := d1; ; c = c.AddDate(0, 0, 1) {
		_, _, rd := c.Date()
		rw := weeks[c.Weekday().String()]
		nextWeeksList = append(nextWeeksList, lbonus.Day{
			Day:               rd,
			DayOfTheWeek:      rw,
			SpecialDay:        false,
			SpecialImageAsset: "",
			Received:          false,
			AdReceived:        false,
			Item: lbonus.Item{
				ItemID:  4,
				AddType: 3001,
				Amount:  1,
			},
		})
		if c.Equal(d2) {
			break
		}
	}

	resp := lbonus.ExecuteResp{
		ResponseData: lbonus.ExecuteData{
			Sheets: []any{},
			CalendarInfo: lbonus.CalendarInfo{
				CurrentDate: time.Now().Format("2006-01-02 03:04:05"),
				CurrentMonth: lbonus.Month{
					Year:  y,
					Month: int(cm),
					Days:  weeksList,
				},
				NextMonth: lbonus.Month{
					Year:  y,
					Month: int(m),
					Days:  nextWeeksList,
				},
			},
			TotalLoginInfo: lbonus.TotalLoginInfo{
				LoginCount:     2626,
				RemainingCount: 74,
				Reward: []lbonus.Reward{
					{
						ItemID:  5,
						AddType: 1000,
						Amount:  5,
					},
				},
			},
			LicenseLbonusList: []any{},
			ClassSystem: lbonus.ClassSystem{
				RankInfo: lbonus.RankInfo{
					BeforeClassRankID: 10,
					AfterClassRankID:  10,
					RankUpDate:        "2020-02-12 11:57:15",
				},
				CompleteFlag: false,
				IsOpened:     true,
				IsVisible:    true,
			},
			StartDashSheets: []any{},
			EffortPoint: []lbonus.EffortPoint{
				{
					LiveEffortPointBoxSpecID: 5,
					Capacity:                 4000000,
					Before:                   1400116,
					After:                    1400116,
					Rewards:                  []lbonus.Rewards{},
				},
			},
			LimitedEffortBox: []any{},
			MuseumInfo:       lbonus.Museum{},
			ServerTimestamp:  time.Now().Unix(),
			PresentCnt:       0,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(resp)
}

func init() {
	router.AddHandler("main.php", "POST", "/lbonus/execute", middleware.Common, execute)
}

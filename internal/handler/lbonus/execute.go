package lbonus

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	lbonusschema "honoka-chan/internal/schema/lbonus"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func execute(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

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

	weeksList := []lbonusschema.Day{}
	for c := d1; ; c = c.AddDate(0, 0, 1) {
		_, _, rd := c.Date()
		received := false
		if rd <= d {
			received = true
		}
		rw := weeks[c.Weekday().String()]
		weeksList = append(weeksList, lbonusschema.Day{
			Day:               rd,
			DayOfTheWeek:      rw,
			SpecialDay:        false,
			SpecialImageAsset: "",
			Received:          received,
			AdReceived:        false,
			Item: lbonusschema.Item{
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

	nextWeeksList := []lbonusschema.Day{}
	for c := d1; ; c = c.AddDate(0, 0, 1) {
		_, _, rd := c.Date()
		rw := weeks[c.Weekday().String()]
		nextWeeksList = append(nextWeeksList, lbonusschema.Day{
			Day:               rd,
			DayOfTheWeek:      rw,
			SpecialDay:        false,
			SpecialImageAsset: "",
			Received:          false,
			AdReceived:        false,
			Item: lbonusschema.Item{
				ItemID:  4,
				AddType: 3001,
				Amount:  1,
			},
		})
		if c.Equal(d2) {
			break
		}
	}

	resp := lbonusschema.ExecuteResp{
		ResponseData: lbonusschema.ExecuteData{
			Sheets: []any{},
			CalendarInfo: lbonusschema.CalendarInfo{
				CurrentDate: time.Now().Format("2006-01-02 03:04:05"),
				CurrentMonth: lbonusschema.Month{
					Year:  y,
					Month: int(cm),
					Days:  weeksList,
				},
				NextMonth: lbonusschema.Month{
					Year:  y,
					Month: int(m),
					Days:  nextWeeksList,
				},
			},
			TotalLoginInfo: lbonusschema.TotalLoginInfo{
				LoginCount:     2626,
				RemainingCount: 74,
				Reward: []lbonusschema.Reward{
					{
						ItemID:  5,
						AddType: 1000,
						Amount:  5,
					},
				},
			},
			LicenseLbonusList: []any{},
			ClassSystem: lbonusschema.ClassSystem{
				RankInfo: lbonusschema.RankInfo{
					BeforeClassRankID: 10,
					AfterClassRankID:  10,
					RankUpDate:        "2020-02-12 11:57:15",
				},
				CompleteFlag: false,
				IsOpened:     true,
				IsVisible:    true,
			},
			StartDashSheets: []any{},
			EffortPoint: []lbonusschema.EffortPoint{
				{
					LiveEffortPointBoxSpecID: 5,
					Capacity:                 4000000,
					Before:                   1400116,
					After:                    1400116,
					Rewards:                  []lbonusschema.Rewards{},
				},
			},
			LimitedEffortBox: []any{},
			MuseumInfo:       lbonusschema.Museum{},
			ServerTimestamp:  time.Now().Unix(),
			PresentCnt:       0,
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	}

	ss.Respond(resp)
}

func init() {
	router.AddHandler("main.php", "POST", "/lbonus/execute", middleware.Common, execute)
}

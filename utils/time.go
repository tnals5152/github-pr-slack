package utils

import "time"

// pr 시간이 지났는지 확인
// 시간이 지났으면 true로 리턴
func IsAfterDay(compareTime time.Time, passDay int) bool {
	now := time.Now().Add(9 * time.Hour)
	PRDay := passDay
	// 리뷰 요청 후 N일 지나면 slack 알림 전송
	for i := 0; i <= passDay; i++ {
		nextDay := compareTime.AddDate(0, 0, i)
		if nextDay.Weekday() == time.Saturday || nextDay.Weekday() == time.Sunday {
			PRDay++
		}
	}

	PRtime := compareTime.AddDate(0, 0, PRDay)

	for {
		if PRtime.Weekday() == time.Saturday || PRtime.Weekday() == time.Sunday {
			PRtime = PRtime.AddDate(0, 0, 1)
			continue
		}
		break
	}

	return PRtime.Before(now)
}

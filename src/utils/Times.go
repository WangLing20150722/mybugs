package utils

import "time"

/**
将精确的时间，转换成mantis List页面的时间显示:
2018-05-20 10:12:30  --> 2018-05-20
 */
func FormatTime2Day(exact time.Time) (time.Time,error) {
	return time.Parse("2006-01-02", exact.Format("2006-01-02"))
}

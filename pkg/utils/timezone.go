package utils

import "time"

var ChinaStandardTimezone = time.FixedZone("CST", 8*3600)

// GetShanghaiTimeLocation 获取上海的
func GetShanghaiTimeLocation() *time.Location {
	return ChinaStandardTimezone
}

func GetTodayBeginSecond() int64 {
	now := GetTimeBeijing()
	begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return begin.Unix()
}

func GetTimeBeijing() time.Time {
	time.Local = GetShanghaiTimeLocation()
	return time.Now().Local()
}

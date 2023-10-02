package utils

import "time"

type Time int64

func GetTimestamp() Time {
	return Time(time.Now().Unix())
}

func FromTimestamp(timestamp int64) Time {
	return Time(timestamp)
}

func (t Time) Passed() bool {
	return time.Now().Unix() > int64(t)
}

func (t Time) Int64() int64 {
	return int64(t)
}

func (t Time) Add(time time.Duration) Time {
	return Time(t.Int64() + time.Milliseconds())
}

func (t Time) Format() string {
	timestamp := time.Now().Unix()
	unixTimestamp := time.Unix(timestamp, 0)
	return unixTimestamp.Format("02.01.2006 | 15:04:05")
}

func (t Time) FormattedDate() string {
	timestamp := time.Now().Unix()
	unixTimestamp := time.Unix(timestamp, 0)
	return unixTimestamp.Format("02.01.2006")
}

func (t Time) FormattedTime() string {
	timestamp := time.Now().Unix()
	unixTimestamp := time.Unix(timestamp, 0)
	return unixTimestamp.Format("15:04:05")
}

func (t Time) Weekday() string {
	timestamp := time.Now().Unix()
	unixTimestamp := time.Unix(timestamp, 0)
	weekday := unixTimestamp.Weekday()
	return weekday.String()
}

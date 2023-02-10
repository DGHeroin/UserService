package utils

import "time"

func GetTime() time.Time {
    return time.Now()
}
func GetUnixTime() int64 {
    return GetTime().Unix()
}

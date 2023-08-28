package util

import "time"

func MsToTime(stamp uint32) time.Time {
	return time.Unix(int64(stamp), 0).Add(946684800 * time.Second)
}

func MsTimeTo(t time.Time) []byte {
	return To4Bytes(t.Add(-946684800 * time.Second).Unix())
}

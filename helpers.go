package uiscom

import (
	"encoding/json"
	"fmt"
	"time"
)

const DateFormat = "2006-01-02 15:04:05"

func TimeToString(t time.Time) string {
	return t.Format(DateFormat)
}

func StringToTime(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

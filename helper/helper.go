package helper

import "encoding/json"

func PrettyLog(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
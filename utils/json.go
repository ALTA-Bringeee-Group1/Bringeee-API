package utils

import "encoding/json"

func JSONEncode(obj interface{}) string {
	json, _ := json.MarshalIndent(obj, "", "  ")
	return string(json)
}

package utils

import "encoding/json"

func JsonEncode(obj interface{}) string {
	json, _ := json.MarshalIndent(obj, "", "  ")
	return string(json)
}

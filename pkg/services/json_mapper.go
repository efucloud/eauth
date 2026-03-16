package services

import jsoniter "github.com/json-iterator/go"

func copyByJSON(src interface{}, dst interface{}) {
	data, err := jsoniter.Marshal(src)
	if err != nil {
		return
	}
	_ = jsoniter.Unmarshal(data, dst)
}

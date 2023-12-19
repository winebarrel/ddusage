package util

import (
	"encoding/json"
)

func RespToMap(resp json.Marshaler) map[string]any {
	js, err := resp.MarshalJSON()

	if err != nil {
		panic(err)
	}

	m := map[string]any{}
	err = json.Unmarshal(js, &m)

	if err != nil {
		panic(err)
	}

	return m
}

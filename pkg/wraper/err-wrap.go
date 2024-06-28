package wraper

import "encoding/json"

func WrapErr(errContext error) []byte {
	answer := map[string]string{
		"id": "",
		"error": errContext.Error(),
	}
	jsonData, _ := json.Marshal(answer)
	return jsonData
}
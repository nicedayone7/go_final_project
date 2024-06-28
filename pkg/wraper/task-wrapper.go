package wraper

import (
	"encoding/json"
	"fmt"
)

func WrapTask(id int) []byte {
	answer := map[string]string{
		"id": fmt.Sprintf("%d", id),
		"error": "",
	}
	jsonData, _ := json.Marshal(answer)
	return jsonData
}
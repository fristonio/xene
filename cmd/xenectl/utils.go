package main

import (
	"encoding/json"
	"fmt"
)

func prettyPrintJSON(jsonStr string) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println(jsonStr)
		return
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(jsonStr)
	} else {
		fmt.Println(string(b))
	}
}

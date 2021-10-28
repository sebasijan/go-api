package main

import (
	"encoding/json"
	"fmt"
)

func PrintObject(input interface{}) {
	fmt.Print(getMarshalIndent(input))
}

func getMarshalIndent(input interface{}) string {
	s, _ := json.MarshalIndent(input, "", "\t")
	return string(s)
}

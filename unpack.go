package main

import(
	"encoding/json"
	"fmt"
)

func unpackJSON(ResponseData string) ContributorResponses{
	var data ContributorResponses
	fmt.Println(ResponseData)
	err := json.Unmarshal([]byte(ResponseData), &data)
	if err != nil{
		println("There was an error in the Unmarshalling phase: %s", err)
	}
	fmt.Println(data)
	return data
}

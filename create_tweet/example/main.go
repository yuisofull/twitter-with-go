package main

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/dynamo"
)

func main() {
	// JSON string that you want to unmarshal
	jsonStr := `{
		"userID": {"S": "user1"},
		"ID": {"S": "id1"}
	}`

	// Variable to hold the unmarshaled dynamo.PagingKey
	var pagingKey dynamo.PagingKey

	// Unmarshal the JSON string into the dynamo.PagingKey
	err := json.Unmarshal([]byte(jsonStr), &pagingKey)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Print the unmarshaled dynamo.PagingKey
	for key, value := range pagingKey {
		fmt.Printf("Key: %s, Value: %s\n", key, *value.S)
	}
}

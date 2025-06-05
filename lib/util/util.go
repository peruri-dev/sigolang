package util

import (
	"encoding/json"
	"fmt"
)

func MarshalAndUnmarshal(input interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("error marshalling input: %w", err)
	}

	if err := json.Unmarshal(jsonData, output); err != nil {
		return fmt.Errorf("error unmarshalling JSON to output struct: %w", err)
	}

	return nil
}

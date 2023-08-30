package utils

import "encoding/json"

func JsonParse[T any](s string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(s), &args); err != nil {
		return *(new(T)), err
	}

	return args, nil
}

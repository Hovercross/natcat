package util

import "encoding/json"

// MustMarshal will either marshal or panic - useful in those times when something can't fail and you still want full unit test coverage
func MustMarshal(input interface{}) []byte {
	out, err := json.Marshal(input)

	if err != nil {
		panic(err)
	}

	return out
}

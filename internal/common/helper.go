package common

import (
	"fmt"
	"strconv"
)

func ParseStringToUint64(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to uint64: %w", err)
	}
	return value, nil
}

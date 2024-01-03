package stockybot

import "fmt"

var (
	ErrInvalidData = fmt.Errorf("invalid data")
	ErrNotFound    = fmt.Errorf("not found")
)

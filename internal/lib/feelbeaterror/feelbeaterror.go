package feelbeaterror

import "fmt"

type FeelBeatError struct {
	DebugMessage string
	UserMessage  ErrorCode
}

func (e *FeelBeatError) Error() string {
	return fmt.Sprintf("%s: %s", e.UserMessage, e.DebugMessage)
}

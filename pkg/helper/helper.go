package helper

import (
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeInput(input string) string {
	p := bluemonday.UGCPolicy() // Create a policy for user-generated content
	return p.Sanitize(input)
}

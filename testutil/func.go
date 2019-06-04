package testutil

import (
	"fmt"
)

// FS returns function signature to compare
func FS(f interface{}) string {
	return fmt.Sprintf("%v", f)
}

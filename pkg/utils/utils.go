package utils

import (
	"github.com/fatih/color"
)

var DebugEnabled bool = false

// StringSliceDel deletes a specific index from a slide of strings
// Taken from https://yourbasic.org/golang/delete-element-slice/
func StringSliceDel(a []string, i int) []string {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
	a[len(a)-1] = ""     // Erase last element (write zero value).
	a = a[:len(a)-1]     // Truncate slice.
	return a
}

func Debug(msg interface{}) {
	if !DebugEnabled {
		return
	}
	color.Magenta("%v", msg)
}

package model

import "fmt"

// Issue represents YouTrack issue.
type Issue struct {
	ID    string
	Title string
}

// FullID returns uniq ID of issue.
// Must depends on all issue fields.
func (i Issue) FullID() string {
	return fmt.Sprintf("%v %v", i.ID, i.Title)
}

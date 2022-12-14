package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("Models: No matching record found!")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

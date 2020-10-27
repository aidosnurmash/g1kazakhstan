package models

import (
	"time"
)

type Picture struct {
	Id        int64
	Path      string
	Origin    string
	CreatedAt time.Time
}


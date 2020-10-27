package models

import "time"

type Part struct {
	Id int64
	PictureId  int64
	PartNum int
	Path  string
	CreatedAt time.Time
}
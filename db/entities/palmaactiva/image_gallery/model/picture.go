//
// Code generated by go-jet DO NOT EDIT.
// Generated at Tuesday, 10-Dec-19 10:58:12 CET
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Picture struct {
	PictureID   int32 `sql:"primary_key"`
	UserID      int32
	ImageID     int32
	Title       string
	Description *string
	Created     *time.Time
	NumLikes    int32
	NumComments int32
}
//
// Code generated by go-jet DO NOT EDIT.
// Generated at Saturday, 14-Dec-19 13:49:37 CET
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

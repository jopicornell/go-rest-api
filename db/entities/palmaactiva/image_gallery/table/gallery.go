//
// Code generated by go-jet DO NOT EDIT.
// Generated at Tuesday, 10-Dec-19 10:58:12 CET
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var Gallery = newGalleryTable()

type GalleryTable struct {
	postgres.Table

	//Columns
	GalleryID postgres.ColumnInteger
	UserID    postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new GalleryTable with assigned alias
func (a *GalleryTable) AS(alias string) *GalleryTable {
	aliasTable := newGalleryTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newGalleryTable() *GalleryTable {
	var (
		GalleryIDColumn = postgres.IntegerColumn("gallery_id")
		UserIDColumn    = postgres.IntegerColumn("user_id")
	)

	return &GalleryTable{
		Table: postgres.NewTable("image_gallery", "gallery", GalleryIDColumn, UserIDColumn),

		//Columns
		GalleryID: GalleryIDColumn,
		UserID:    UserIDColumn,

		AllColumns:     postgres.ColumnList{GalleryIDColumn, UserIDColumn},
		MutableColumns: postgres.ColumnList{UserIDColumn},
	}
}

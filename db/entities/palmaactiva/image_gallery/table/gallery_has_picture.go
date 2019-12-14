//
// Code generated by go-jet DO NOT EDIT.
// Generated at Thursday, 12-Dec-19 11:57:22 CET
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var GalleryHasPicture = newGalleryHasPictureTable()

type GalleryHasPictureTable struct {
	postgres.Table

	//Columns
	GalleryID postgres.ColumnInteger
	PictureID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new GalleryHasPictureTable with assigned alias
func (a *GalleryHasPictureTable) AS(alias string) *GalleryHasPictureTable {
	aliasTable := newGalleryHasPictureTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newGalleryHasPictureTable() *GalleryHasPictureTable {
	var (
		GalleryIDColumn = postgres.IntegerColumn("gallery_id")
		PictureIDColumn = postgres.IntegerColumn("picture_id")
	)

	return &GalleryHasPictureTable{
		Table: postgres.NewTable("image_gallery", "gallery_has_picture", GalleryIDColumn, PictureIDColumn),

		//Columns
		GalleryID: GalleryIDColumn,
		PictureID: PictureIDColumn,

		AllColumns:     postgres.ColumnList{GalleryIDColumn, PictureIDColumn},
		MutableColumns: postgres.ColumnList{GalleryIDColumn, PictureIDColumn},
	}
}

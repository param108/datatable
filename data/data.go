package data

import (
	"context"
)

type ColumnType string

const (
	//ColumnTypeString - String value
	ColumnTypeString = ColumnType("string")

	//ColumnTypeFloat - Float value
	ColumnTypeFloat = ColumnType("float")

	//ColumnTypeInt - Integer value
	ColumnTypeInt = ColumnType("integer")
)

type ColumnMetadata struct {
	Name string
	Type ColumnType
}

type DataSource interface {
	Create(ctx context.Context) error
	Get(row, col int) (interface{}, error)
	GetColumns() []*ColumnMetadata
	SetColumns(columns []*ColumnMetadata)
	SetColumn(key string, t ColumnType) error
	GetSize() (numRows int, numCols int)
	GetColumn(col int) ([]*Metadata, error)
	GetRow(row int) ([]*Metadata, error)
	Changed() bool
	Set(row, col int, value interface{}) error
}

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
	CreateContext() context.Context
	Create(ctx context.Context) error
	Get(ctx context.Context, row, col int) (interface{}, error)
	GetColumns(ctx context.Context) []*ColumnMetadata
	SetColumns(ctx context.Context, columns []*ColumnMetadata)
	SetColumn(ctx context.Context, key string, t ColumnType) error
}

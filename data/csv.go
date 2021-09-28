package data

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/param108/datatable/types"
	"github.com/pkg/errors"
)

type CSV struct {
	setup    bool
	Filename string
	data     [][]string
	metadata []*ColumnMetadata
}

func NewCSV(filename string) (DataSource, error) {
	c := &CSV{}
	ctx := c.CreateContext()
	ctx = context.WithValue(ctx, types.ContextKey("filename"), filename)
	if err := c.Create(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

//CreateContext - returns the context to be used
// for a CSV datasource
func (c *CSV) CreateContext() context.Context {
	return context.Background()
}

//Create - Instantiates the csv data source
// Required Keys in Context
// filename - the filename for the data
func (c *CSV) Create(ctx context.Context) error {
	filename, ok := ctx.Value(types.ContextKey("filename")).(string)
	if !ok {
		return errors.New("invalid filename")
	}

	fp, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "open csv")
	}
	r := csv.NewReader(fp)
	data, err := r.ReadAll()
	if err != nil {
		return errors.Wrap(err, "read csv")
	}

	c.data = data

	c.createMetadata()

	c.setup = true
	return nil
}

// Get - row and col are zero indexed
func (c *CSV) Get(ctx context.Context, row, col int) (interface{}, error) {
	if len(c.data) <= row {
		return nil, errors.New("row idx too large")
	}

	if len(c.data[row]) <= col {
		return nil, errors.New("col idx too large")
	}

	// This is the header row
	if row == 0 {
		return nil, errors.New("invalid row")
	}
	return c.data[row][col], nil
}

func (c *CSV) createMetadata() {
	c.metadata = []*ColumnMetadata{}
	for _, header := range c.data[0] {
		cm := &ColumnMetadata{}
		cm.Name = header
		cm.Type = ColumnTypeString
		c.metadata = append(c.metadata, cm)
	}
}

// GetColumns - get list of columns and their types
func (c *CSV) GetColumns(ctx context.Context) []*ColumnMetadata {
	return c.metadata
}

// SetColumns - set the types of all the columns in one shot
func (c *CSV) SetColumns(ctx context.Context, columns []*ColumnMetadata) {
	c.metadata = columns
}

// SetColumn - set the Column type for one column
func (c *CSV) SetColumn(ctx context.Context, key string, t ColumnType) error {
	for _, m := range c.metadata {
		if m.Name == key {
			m.Type = t
			return nil
		}
	}

	return errors.New("invalid key")
}

// GetSize - returns the size of the data table
func (c *CSV) GetSize(ctx context.Context) (numRows, numCols int) {
	numRows = len(c.data)
	if len(c.data) == 0 {
		return numRows, numCols
	}
	numCols = len(c.data[0])
	return numRows, numCols
}

//GetColumn - returns the data for a column
func (c *CSV) GetColumn(ctx context.Context, col int) ([]string, error) {
	if len(c.data) == 0 {
		return nil, errors.New("invalid column")
	}
	if len(c.data[0]) <= col {
		return nil, errors.New("invalid column")
	}

	ret := []string{}

	for _, row := range c.data {
		ret = append(ret, row[col])
	}
	return ret, nil
}

//GetRow - returns the data for a Row
func (c *CSV) GetRow(ctx context.Context, row int) ([]string, error) {
	if len(c.data) == 0 {
		return nil, errors.New("invalid column")
	}

	return c.data[row], nil
}

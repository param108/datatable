package data

import (
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
	"sync"

	"github.com/param108/datatable/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type CSV struct {
	setup          bool
	Filename       string
	data           [][]string
	Data           [][]*Metadata
	columnMetadata []*ColumnMetadata
	changed        bool
	mx             sync.RWMutex
}

type Metadata struct {
	Value   string
	FgColor string
	BgColor string
	Attr    string
}

func NewCSV(filename string) (DataSource, error) {
	c := &CSV{Filename: filename}
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ContextKey("filename"), filename)
	if err := c.Create(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CSV) Save() error {
	file, err := ioutil.TempFile("", "datatable")
	if err != nil {
		log.Errorf("failed to open temp file %v", err)
		return err
	}

	oldpath := file.Name()
	writer := csv.NewWriter(file)

	err = writer.WriteAll(c.data)
	if err != nil {
		log.Errorf("failed to write file %v", err)
		return err
	}

	file.Close()

	os.Rename(oldpath, c.Filename)
	return nil
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

	c.Data = [][]*Metadata{}
	for rowNum, rowData := range c.data {
		rowM := []*Metadata{}
		for _, d := range rowData {
			m := &Metadata{
				Value: d,
			}

			if rowNum == 0 {
				m.BgColor = "[_light_gray_]"
				m.FgColor = "[black]"
				m.Attr = "[underline]"
			}
			rowM = append(rowM, m)
		}
		c.Data = append(c.Data, rowM)
	}

	c.createMetadata()

	c.setup = true
	return nil
}

func (c *CSV) Set(row, col int, value interface{}) error {
	if len(c.data) <= row {
		return errors.New("row idx too large")
	}

	if len(c.data[row]) <= col {
		return errors.New("col idx too large")
	}

	// This is the header row
	if row == 0 {
		return errors.New("invalid row")
	}

	if v, ok := value.(string); ok {
		c.mx.Lock()
		c.data[row][col] = v
		c.Data[row][col].Value = v
		c.changed = true
		c.mx.Unlock()
	} else {
		return errors.New("invalid value")
	}

	return nil
}

// Get - row and col are zero indexed
func (c *CSV) Get(row, col int) (interface{}, error) {
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

	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.data[row][col], nil
}

func (c *CSV) createMetadata() {
	c.columnMetadata = []*ColumnMetadata{}
	for _, header := range c.data[0] {
		cm := &ColumnMetadata{}
		cm.Name = header
		cm.Type = ColumnTypeString
		c.columnMetadata = append(c.columnMetadata, cm)
	}
}

// GetColumns - get list of columns and their types
func (c *CSV) GetColumns() []*ColumnMetadata {
	return c.columnMetadata
}

// SetColumns - set the types of all the columns in one shot
func (c *CSV) SetColumns(columns []*ColumnMetadata) {
	c.columnMetadata = columns
}

// SetColumn - set the Column type for one column
func (c *CSV) SetColumn(key string, t ColumnType) error {
	for _, m := range c.columnMetadata {
		if m.Name == key {
			m.Type = t
			return nil
		}
	}

	return errors.New("invalid key")
}

// GetSize - returns the size of the data table
func (c *CSV) GetSize() (numRows, numCols int) {
	numRows = len(c.data)
	if len(c.data) == 0 {
		return numRows, numCols
	}
	numCols = len(c.data[0])
	return numRows, numCols
}

//GetColumn - returns the data for a column
func (c *CSV) GetColumn(col int) ([]*Metadata, error) {
	if len(c.data) == 0 {
		return nil, errors.New("invalid column")
	}
	if len(c.data[0]) <= col {
		return nil, errors.New("invalid column")
	}

	ret := []*Metadata{}

	for _, row := range c.Data {
		ret = append(ret, row[col])
	}
	return ret, nil
}

//GetRow - returns the data for a Row
func (c *CSV) GetRow(row int) ([]*Metadata, error) {
	if len(c.data) == 0 {
		return nil, errors.New("invalid column")
	}

	return c.Data[row], nil
}

func (c *CSV) Changed() bool {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.changed
}

func (c *CSV) ClearChanged() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.changed = false
}

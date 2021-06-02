package rqlite

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
)

// QueryResult hold any errors from that query, a list of columns and types, and the actual row values.
type QueryResult struct {
	Err       error
	columns   []string
	types     []string
	Timing    float64
	values    []interface{}
	rowNumber int64
}

// Columns returns a list of the column names for this QueryResult.
func (qr *QueryResult) Columns() []string {
	return qr.columns
}

// Map() returns the current row (as advanced by Next()) as a map[string]interface{}
// The key is a string corresponding to a column name.
// The value is the corresponding column.
func (qr *QueryResult) Map() (map[string]interface{}, error) {
	ans := make(map[string]interface{})

	if qr.rowNumber == -1 {
		return ans, errors.New("need to Next() before you Map()")
	}

	thisRowValues := qr.values[qr.rowNumber].([]interface{})
	for i := 0; i < len(qr.columns); i++ {
		switch qr.types[i] {
		case "date", "datetime":
			if thisRowValues[i] != nil {
				t, err := toTime(thisRowValues[i])
				if err != nil {
					return ans, err
				}
				ans[qr.columns[i]] = t
			} else {
				ans[qr.columns[i]] = nil
			}
		default:
			ans[qr.columns[i]] = thisRowValues[i]
		}
	}

	return ans, nil
}

func toTime(src interface{}) (time.Time, error) {
	switch src := src.(type) {
	case string:
		const layout = "2006-01-02 15:04:05"
		if t, err := time.Parse(layout, src); err == nil {
			return t, nil
		}
		return time.Parse(time.RFC3339, src)
	case float64:
		return time.Unix(int64(src), 0), nil
	case int64:
		return time.Unix(src, 0), nil
	}
	return time.Time{}, fmt.Errorf("invalid time type:%T val:%v", src, src)
}

// Next() positions the QueryResult result pointer so that Scan() or Map() is ready.
// You should call Next() first, but gorqlite will fix it if you call Map() or Scan() before
// the initial Next().
func (qr *QueryResult) Next() bool {
	if qr.rowNumber >= int64(len(qr.values)-1) {
		return false
	}

	qr.rowNumber += 1
	return true
}

// NumRows() returns the number of rows returned by the query.
func (qr *QueryResult) NumRows() int64 {
	return int64(len(qr.values))
}

// RowNumber() returns the current row number as Next() iterates through the result's rows.
func (qr *QueryResult) RowNumber() int64 {
	return qr.rowNumber
}

// Scan() takes a list of pointers and then updates them to reflect he current row's data.

// Note that only the following data types are used, and they
// are a subset of the types JSON uses:
// 	string, for JSON strings
// 	float64, for JSON numbers
// 	int64, as a convenient extension
// 	nil for JSON null

// booleans, JSON arrays, and JSON objects are not supported,
// since sqlite does not support them.
func (qr *QueryResult) Scan(dest ...interface{}) error {
	if qr.rowNumber == -1 {
		return errors.New("need to Next() before Scan()")
	}

	if len(dest) != len(qr.columns) {
		return errors.Errorf("expected %d columns but got %d vars", len(qr.columns), len(dest))
	}

	thisRowValues := qr.values[qr.rowNumber].([]interface{})
	for n, d := range dest {
		src := thisRowValues[n]
		if src == nil {
			continue
		}
		switch d := d.(type) {
		case *time.Time:
			if src == nil {
				continue
			}
			t, err := toTime(src)
			if err != nil {
				return errors.Errorf("%v: bad time col:(%d/%s) val:%v", err, n, qr.Columns()[n], src)
			}
			*d = t
		case *int:
			switch src := src.(type) {
			case float64:
				*d = int(src)
			case int64:
				*d = int(src)
			case string:
				i, err := strconv.Atoi(src)
				if err != nil {
					return err
				}
				*d = i
			default:
				return errors.Errorf("invalid int col:%d type:%T val:%v", n, src, src)
			}
		case *int64:
			switch src := src.(type) {
			case float64:
				*d = int64(src)
			case int64:
				*d = src
			case string:
				i, err := strconv.ParseInt(src, 10, 64)
				if err != nil {
					return err
				}
				*d = i
			default:
				return errors.Errorf("invalid int64 col:%d type:%T val:%v", n, src, src)
			}
		case *float64:
			*d = float64(src.(float64))
		case *string:
			switch src := src.(type) {
			case string:
				*d = src
			default:
				return errors.Errorf("invalid string col:%d type:%T val:%v", n, src, src)
			}
		default:
			return errors.Errorf("unknown destination type (%T) to scan into in variable #%d", d, n)
		}
	}

	return nil
}

// Types() returns an array of the column's types.
func (qr *QueryResult) Types() []string {
	return qr.types
}

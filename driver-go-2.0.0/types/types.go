package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/taosdata/driver-go/v2/errors"
	"time"
)

type NullInt64 struct {
	Inner int64
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(int64)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse int64 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullInt32 struct {
	Inner int32
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt32) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(int32)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse int32 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullInt16 struct {
	Inner int16
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt16) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(int16)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse int16 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullInt8 struct {
	Inner int8
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt8) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(int8)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse int8 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullUInt64 struct {
	Inner uint64
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUInt64) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(uint64)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse uint64 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullUInt32 struct {
	Inner uint32
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUInt32) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(uint32)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse uint32 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullUInt16 struct {
	Inner uint16
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUInt16) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(uint16)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse uint16 error"}
	}
	n.Inner = v
	return nil
}

// Value implements the driver Valuer interface.
func (n NullUInt16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullUInt8 struct {
	Inner uint8
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullUInt8) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(uint8)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse uint8 error"}
	}
	n.Inner = v
	return nil
}

type NullFloat32 struct {
	Inner float32
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(float32)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse float32 error"}
	}
	n.Inner = v
	return nil
}

func (n NullFloat32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullFloat64 struct {
	Inner float64
	Valid bool // Valid is true if Inner is not NULL
}

// Scan implements the Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Inner, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	v, ok := value.(float64)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse float32 error"}
	}
	n.Inner = v
	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullBool struct {
	Inner bool
	Valid bool // Valid is true if Inner is not NULL
}

func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	v, ok := value.(bool)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse float32 error"}
	}
	n.Inner = v
	return nil
}

func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

type NullString struct {
	Inner string
	Valid bool // Valid is true if Inner is not NULL
}

func (n *NullString) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	n.Valid = true
	v, ok := value.(string)
	if !ok {
		return &errors.TaosError{Code: 0xffff, ErrStr: "taosSql parse float32 error"}
	}
	n.Inner = v
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

// Value implements the driver Valuer interface.
func (n NullUInt8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Inner, nil
}

func (n NullUInt8) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}
func (n NullUInt16) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullUInt32) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullUInt64) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}
func (n NullInt8) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}
func (n NullInt16) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullInt32) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullInt64) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullFloat32) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

func (n NullFloat64) String() string {
	if n.Valid {
		return fmt.Sprintf("%v", n.Inner)
	}
	return "NULL"
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
// The value type must be time.Time or string / []byte (formatted time-string),
// otherwise Scan fails.
func (nt *NullTime) Scan(value interface{}) (err error) {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = v, true
		return
	case []byte:
		nt.Time, err = time.Parse(time.RFC3339Nano, string(v))
		nt.Valid = err == nil
		return
	case string:
		nt.Time, err = time.Parse(time.RFC3339Nano, v)
		nt.Valid = err == nil
		return
	}

	nt.Valid = false
	return fmt.Errorf("can't convert %T to time.Time", value)
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

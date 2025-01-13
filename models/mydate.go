package models

import (
    "database/sql/driver"
    "time"
    "fmt"
)

// Custom Date type that implements the Scanner interface
type MyDate struct {
    time.Time
}

// Implement the Scan method for the custom type
func (d *MyDate) Scan(value interface{}) error {
    if value == nil {
        return nil
    }

    // Check if the value is a byte slice (which is how MySQL typically returns DATE values)
    switch v := value.(type) {
    case []byte:
        // Parse the DATE value (which is a string) into time.Time
        parsedDate, err := time.Parse("2006-01-02", string(v))
        if err != nil {
            return fmt.Errorf("error parsing date: %v", err)
        }
        d.Time = parsedDate
        return nil
    default:
        return fmt.Errorf("unsupported value type for Date field: %T", v)
    }
}

// Implement the Value method for the Valuer interface
func (d MyDate) Value() (driver.Value, error) {
    // Return the time in MySQL DATE format (YYYY-MM-DD)
    return d.Time.Format("2006-01-02"), nil
}
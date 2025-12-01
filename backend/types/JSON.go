package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for JSONB")
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) GetString(key string) (string, bool) {
	if j == nil {
		return "", false
	}
	val, ok := (*j)[key]
	if !ok {
		return "", false
	}
	s, ok := val.(string)
	return s, ok
}

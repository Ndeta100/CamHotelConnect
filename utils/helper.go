package utils

import (
	"strings"
	"time"
)

type JSONTime time.Time

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*jt = JSONTime(t)
	return nil
}

func (jt JSONTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(jt).Format("2006-01-02") + "\""), nil
}

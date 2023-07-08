package models

import (
	"encoding/json"
	"errors"
	"time"
)

type NewTime time.Time

func (n *NewTime) UnmarshalJSON(data []byte) (err error) {

	if len(data) == 0 || string(data) == "null" {
		return
	}
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	// location 사용 안 할 때는
	t, err := time.Parse("2006-01-02T15:04:05Z", s)

	if err != nil {
		return
	}

	*n = NewTime(t)
	return
}

func (n NewTime) MarshalJSON() ([]byte, error) {
	t := time.Time(n)
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len("2006-01-02T15:04:05Z")+2)
	b = append(b, '"')
	b = t.AppendFormat(b, "2006-01-02T15:04:05Z")
	b = append(b, '"')
	return b, nil
}

package org

import (
	"time"
)

var (
	timestampFormat = "2006-01-02 Mon 15:04"
	datestampFormat = "2006-01-02 Mon"
)

type Timestamp struct {
	Time     time.Time
	IsDate   bool
	Interval string
}

func (t Timestamp) String() string {
	if t.IsDate {
		return t.Time.Format(datestampFormat)
	}
	return t.Time.Format(timestampFormat)
}

func ParseTimestamp(value, interval string) (Timestamp, error) {
	t, err := time.Parse(timestampFormat, value)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{
		Time:     t,
		IsDate:   false,
		Interval: interval,
	}, nil
}

func ParseDatestamp(value, interval string) (Timestamp, error) {
	t, err := time.Parse(datestampFormat, value)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{
		Time:     t,
		IsDate:   true,
		Interval: interval,
	}, nil
}

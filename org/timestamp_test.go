package org_test

import (
	"testing"

	"github.com/Ladicle/org2html/org"
)

func mustParseTimestamp(t *testing.T, val, interval string) org.Timestamp {
	tp, err := org.ParseTimestamp(val, interval)
	if err != nil {
		t.Fatalf("fail to ParseTimestamp(): err=%v", err)
	}
	return tp
}

func mustParseDatestamp(t *testing.T, val, interval string) org.Timestamp {
	tp, err := org.ParseDatestamp(val, interval)
	if err != nil {
		t.Fatalf("fail to ParseTimestamp(): err=%v", err)
	}
	return tp
}

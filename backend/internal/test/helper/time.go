package helper

import (
	"testing"
	"time"
)

func Time(t *testing.T, value string) time.Time {
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04", value, time.UTC)
	if err != nil {
		t.Fatalf("時間のパースに失敗しました: %v", err)
	}
	return parsedTime.UTC()
}

func TimePtr(t *testing.T, value string) *time.Time {
	parsedTime := Time(t, value)
	return &parsedTime
}

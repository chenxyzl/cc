package chat

import "testing"

func TestFormatTime(t *testing.T) {
	if formatOutTime(0) != "00d 00h 00m 00s" {
		t.Fail()
	}
	if formatOutTime(100) != "00d 00h 01m 40s" {
		t.Fail()
	}
	if formatOutTime(60) != "00d 00h 01m 00s" {
		t.Fail()
	}
	if formatOutTime(600) != "00d 00h 10m 00s" {
		t.Fail()
	}
	if formatOutTime(3600) != "00d 01h 00m 00s" {
		t.Fail()
	}
	if formatOutTime(36061) != "00d 10h 01m 01s" {
		t.Fail()
	}
}

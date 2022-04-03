package main

import (
	"testing"
)

func Test_allowedChar(t *testing.T) {
	//try a unit test on allowdChar
	if c := allowedChar(); c < 48 || (57 < c && c < 65) || (90 < c && c < 97) {
		t.Error("allowedChar() return not in [0-9A-Za-z]")
	} else {
		t.Log("allowedChar() success with:", c)
	}
}

func Test_randomString(t *testing.T) {
	//try a unit test on randomString
	Length := 6
	if s := randomString(Length); len(s) != Length {
		t.Error("randomString() return error string")
	} else {
		t.Log("randomString() success return string:", s)
	}
}

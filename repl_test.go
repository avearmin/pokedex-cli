package main

import (
	"testing"
)

func TestTrimAndLower(t *testing.T) {
	cleanString := "gopher"
	dirtyString := "  GoPher    "
	trimedAndLowered := trimAndLower(dirtyString)
	if cleanString != trimAndLower(dirtyString) {
		t.Errorf("Expected: %s, but got: %s", cleanString, trimedAndLowered)
	}
	cleanString = "go is awesome"
	dirtyString = "                   Go IS AWesOME           "
	trimedAndLowered = trimAndLower(dirtyString)
	if cleanString != trimAndLower(dirtyString) {
		t.Errorf("Expected: %s, but got: %s", cleanString, trimedAndLowered)
	}
}

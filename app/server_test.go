package main

import "testing"

func TestParseInt(t *testing.T) {
	ty, d, err := parse(":23")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respInt {
		t.Errorf("Incorrect type")
	}

	if d != 23 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseIntError(t *testing.T) {
	_, _, err := parse(":1a")

	if err == nil {
		t.Errorf("Unexpected error")
	}

}

func TestParseIntTrim(t *testing.T) {
	ty, d, err := parse(" :23 ")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respInt {
		t.Errorf("Incorrect type")
	}

	if d != 23 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseString(t *testing.T) {
	ty, d, err := parse("+Hotel")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respString {
		t.Errorf("Incorrect type")
	}

	if d != "Hotel" {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseArray(t *testing.T) {
	ty, d, err := parse("*5")

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if ty != respArray {
		t.Errorf("Incorrect type")
	}

	if d != 5 {
		t.Errorf("Incorrect value %v", d)
	}
}

func TestParseArrayError(t *testing.T) {
	_, _, err := parse("*1a")

	if err == nil {
		t.Errorf("Unexpected error")
	}

}

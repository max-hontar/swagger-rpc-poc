package main_test

import (
	"swagger-rpc-poc"
	"testing"
)

func TestCheckInTest(t *testing.T) {
	err := main.CheckInTest(false)
	if err != nil {
		t.FailNow()
	}

	err = main.CheckInTest(true)
	if err == nil {
		t.FailNow()
	}
}

func TestAnotherCheckInTest(t *testing.T) {
	err := main.AnotherCheckInTest(false)
	if err != nil {
		t.FailNow()
	}

	err = main.AnotherCheckInTest(true)
	if err == nil {
		t.FailNow()
	}
}

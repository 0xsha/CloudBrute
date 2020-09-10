package internal

import (
	"testing"
)

func TestCheckSupportedCloud(t *testing.T) {

	config := InitConfig("../config/config.yaml")
	got, err := CheckSupportedCloud("amazon", config)

	if err != nil {

		t.Errorf("Err %s", err)
	}

	want := "amazon"
	if got != want {

		t.Errorf("Err")
	}

}

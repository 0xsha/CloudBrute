package internal

import (
	"github.com/modern-go/reflect2"
	"testing"
)

func TestInitConfig(t *testing.T) {

	want := new(Config)
	got := InitConfig("../config/config.yaml")

	if reflect2.TypeOf(want) != reflect2.TypeOf(got) {

		t.Errorf("Got wrong config format")
	}

}

func TestInitCloudConfig(t *testing.T) {

	config := InitConfig("../config/config.yaml")

	for _, provider := range config.Providers {

		_, err := InitCloudConfig(provider, "../config/modules/")

		if err != nil {
			t.Errorf("Cant load config for " + provider)
		}
	}

}

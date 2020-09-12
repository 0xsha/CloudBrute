package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateMutatedUrls(t *testing.T) {

	envs := []string{"test", "dev", "prod", "stage"}

	got, err := GenerateMutatedUrls("../data/storage_small.txt", "storage", "amazon", "../config/modules/", "target", envs)
	if err != nil {

		t.Errorf("Error generating urls %s", err)
	}

	var stringArr []string

	if reflect.TypeOf(got) != reflect.TypeOf(stringArr) {

		fmt.Println("Received wrong type")
	}

}

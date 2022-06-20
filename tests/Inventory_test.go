package tests

import (
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"testing"
)

func TestYAMLInventory(t *testing.T) {
	inv, err := inventory.LoadFromYAML("secrets/hosts.yml")
	if err != nil {
		panic(err)
	}
	for _, host := range inv.Hosts {
		t.Logf("inventory contains host definition: %s", host)
	}
}

func TestCSVInventory(t *testing.T) {
	inv, err := inventory.LoadFromCSV("secrets/hosts.csv")
	if err != nil {
		panic(err)
	}
	for _, host := range inv.Hosts {
		t.Logf("inventory contains host definition: %s", host)
	}
}

package main

import (
	"testing"

	c "delivery/config"
)

func TestReadOrders(t *testing.T) {
	orders, err := readOrders("orders_test.json")
	if err != nil {
		t.Fatal(err)
	}

	want := 132
	got := len(orders)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCreateShelvesFromConfig(t *testing.T) {
	c.Config.Shelves = []struct {
		Name          string `yaml:"name"`
		Temperature   string `yaml:"temp"`
		Capacity      int    `yaml:"cap"`
		DecayModifier int    `yaml:"decayModifier​"`
	}{{
		Name:          "Frozen shelf",
		Temperature:   "frozen",
		Capacity:      10,
		DecayModifier: 1,
	}, {
		Name:          "Cold shelf",
		Temperature:   "cold",
		Capacity:      10,
		DecayModifier: 1,
	}, {
		Name:          "Hot shelf",
		Temperature:   "hot",
		Capacity:      10,
		DecayModifier: 1,
	}}

	want := len(c.Config.Shelves)
	got := len(createShelvesFromConfig())

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

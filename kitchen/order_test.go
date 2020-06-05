package kitchen

import (
	"testing"
)

func TestOrder(t *testing.T) {
	t.Run("IncAge", func(t *testing.T) {
		order := &Order{
			ID:          "a8cfcb76-7f24-4420-a5ba-d46dd77bdffd",
			Name:        "Banana Split",
			Temperature: "frozen",
			ShelfLife:   20,
			DecayRate:   0.63,
		}

		want := order.age + 1

		order.IncAge()

		got := order.age

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("GetInherentValue", func(t *testing.T) {
		order := &Order{
			ID:          "a8cfcb76-7f24-4420-a5ba-d46dd77bdffd",
			Name:        "Banana Split",
			Temperature: "frozen",
			ShelfLife:   20,
			DecayRate:   0.63,

			shelfDecayModifier: 1,
		}

		order.IncAge()
		order.IncAge()
		order.IncAge()

		want := 0.9055
		got := order.GetInherentValue()

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

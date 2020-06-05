package kitchen

import (
	"testing"
	"time"

	c "delivery/config"
)

func init() {
	c.Init("config.yml")
}

func TestKitchen(t *testing.T) {
	t.Run("PlaceOrder", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 10, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 10, 1)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			NewShelf("Overflow shelf", "any", 15, 3),
		)

		order := &Order{
			ID:          "a8cfcb76-7f24-4420-a5ba-d46dd77bdffd",
			Name:        "Banana Split",
			Temperature: "frozen",
			ShelfLife:   20,
			DecayRate:   0.63,
		}

		want := true
		got := k.PlaceOrder(order)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		} else {
			placed := false
			for orderID := range frozenShelf.orders {
				if orderID == order.ID {
					placed = true
					break
				}
			}

			if !placed {
				t.Errorf("got %v want %v", placed, true)
			}
		}
	})

	t.Run("PlaceOrder_Negative", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 10, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 10, 1)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			NewShelf("Overflow shelf", "any", 15, 3),
		)

		order := &Order{
			ID:          "a8cfcb76-7f24-4420-a5ba-d46dd77bdffd",
			Name:        "Banana Split",
			Temperature: "cold",
			ShelfLife:   20,
			DecayRate:   0.63,
		}

		want := false
		got := k.PlaceOrder(order)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("GetAvailableShelves", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 1, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 10, 1)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			NewShelf("Overflow shelf", "any", 15, 3),
		)

		want := 2
		got := len(k.GetAvailableShelves())
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		frozenShelf.orders["1"] = &Order{}

		want = 1
		got = len(k.GetAvailableShelves())
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 1, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 10, 1)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			NewShelf("Overflow shelf", "any", 15, 3),
		)

		got := k.IsEmpty()
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("IsEmpty_Negative", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 1, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 10, 1)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			NewShelf("Overflow shelf", "any", 15, 3),
		)

		frozenShelf.orders["1"] = &Order{ID: "1", Temperature: "frozen"}

		got := k.IsEmpty()
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("RotateOrdersFromOverflowShelve", func(t *testing.T) {
		frozenShelf := NewShelf("Frozen shelf", "frozen", 1, 1)
		hotShelf := NewShelf("Hot shelf", "hot", 1, 1)
		overflowShelf := NewShelf("Overflow shelf", "any", 2, 2)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
				"hot":    hotShelf,
			},
			overflowShelf,
		)

		frozenShelf.orders["1"] = &Order{ID: "1", Temperature: "frozen"}
		overflowShelf.orders["2"] = &Order{ID: "2", Temperature: "hot"}
		overflowShelf.orders["3"] = &Order{ID: "3", Temperature: "hot"}

		order := &Order{
			ID:          "4",
			Temperature: "frozen",
		}
		ok := k.RotateOrdersFromOverflowShelve(order)
		if !ok {
			t.Errorf("got %v want %v", ok, true)
			return
		}

		want := 1
		got := len(hotShelf.orders)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		order = &Order{
			ID:          "5",
			Temperature: "frozen",
		}
		ok = k.RotateOrdersFromOverflowShelve(order)
		if !ok {
			t.Errorf("got %v want %v", ok, true)
			return
		}

		if _, ok := overflowShelf.orders[order.ID]; !ok {
			t.Errorf("got %v want %v", ok, true)
		}
	})

	t.Run("CreateCourier", func(t *testing.T) {
		c.Config.Courier.Arrive.Duration = time.Second
		c.Config.Courier.Arrive.Min = 1
		c.Config.Courier.Arrive.Max = 2

		frozenShelf := NewShelf("Frozen shelf", "frozen", 10, 1)
		overflowShelf := NewShelf("Overflow shelf", "any", 2, 2)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
			},
			overflowShelf,
		)

		order := &Order{
			ID:          "1",
			Temperature: "frozen",
		}
		frozenShelf.orders[order.ID] = order

		k.CreateCourier(order)

		time.Sleep(time.Duration(c.Config.Courier.Arrive.Max) * time.Second)

		if _, ok := frozenShelf.orders[order.ID]; ok {
			t.Errorf("got %v want %v", ok, false)
		}
	})

	t.Run("PickUpOrder", func(t *testing.T) {
		c.Config.Courier.Arrive.Duration = time.Second
		c.Config.Courier.Arrive.Min = 1
		c.Config.Courier.Arrive.Max = 2

		frozenShelf := NewShelf("Frozen shelf", "frozen", 10, 1)
		overflowShelf := NewShelf("Overflow shelf", "any", 2, 2)

		k := New(
			map[string]*Shelf{
				"frozen": frozenShelf,
			},
			overflowShelf,
		)

		order := &Order{
			ID:          "1",
			Temperature: "frozen",
		}
		frozenShelf.orders[order.ID] = order

		ok := k.PickUpOrder(order)
		if !ok {
			t.Errorf("got %v want %v", ok, true)
		}

		if _, ok := frozenShelf.orders[order.ID]; ok {
			t.Errorf("got %v want %v", ok, false)
		}
	})
}

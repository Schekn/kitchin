package kitchen

import (
	"fmt"
	"testing"
)

func TestShelf(t *testing.T) {
	t.Run("OrdersCount", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 10, 1)

		want := 5
		for i := 0; i < want; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}
		}

		got := shelf.OrdersCount()

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("HasEmptySeats", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 10, 1)

		for i := 0; i < 5; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}
		}

		got := shelf.HasEmptySeats()
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("HasEmptySeats_Negative", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		for i := 0; i < 5; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}

		}

		got := shelf.HasEmptySeats()
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		got := shelf.IsEmpty()
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("IsEmpty_Negative", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		for i := 0; i < 5; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}

		}

		got := shelf.IsEmpty()
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("AddOrder", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		for i := 0; i < 3; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}
		}

		got := shelf.AddOrder(&Order{})
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("AddOrder_Negative", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		for i := 0; i < 5; i++ {
			shelf.orders[fmt.Sprintf("%d", i)] = &Order{}
		}

		got := shelf.AddOrder(&Order{})
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("WithdrawOrder", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		ordersCount := 5
		for i := 0; i < ordersCount; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID}
		}

		orderID := "3"
		order, ok := shelf.WithdrawOrder(orderID)

		if order != nil && order.ID != orderID {
			t.Errorf("got %v want %v", order.ID, orderID)
		}

		if !ok {
			t.Errorf("got %v want %v", ok, true)
		}

		if shelf.OrdersCount() != ordersCount-1 {
			t.Errorf("got %v want %v", shelf.OrdersCount(), ordersCount-1)
		}
	})

	t.Run("WithdrawOrder_Negative", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		ordersCount := 6
		for i := 0; i < ordersCount; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID}
		}

		_, ok := shelf.WithdrawOrder("6")

		if ok {
			t.Errorf("got %v want %v", ok, false)
		}

		if shelf.OrdersCount() != ordersCount {
			t.Errorf("got %v want %v", shelf.OrdersCount(), ordersCount)
		}
	})

	t.Run("DeleteOrder", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		ordersCount := 6
		for i := 0; i < ordersCount; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID}
		}

		ok := shelf.DeleteOrder("3")

		if !ok {
			t.Errorf("got %v want %v", ok, true)
		}

		if shelf.OrdersCount() != ordersCount-1 {
			t.Errorf("got %v want %v", shelf.OrdersCount(), ordersCount-1)
		}
	})

	t.Run("DeleteOrder_Negative", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		ordersCount := 6
		for i := 0; i < ordersCount; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID}
		}

		ok := shelf.DeleteOrder("6")

		if ok {
			t.Errorf("got %v want %v", ok, false)
		}

		if shelf.OrdersCount() != ordersCount {
			t.Errorf("got %v want %v", shelf.OrdersCount(), ordersCount)
		}
	})

	t.Run("DeleteRandomOrder", func(t *testing.T) {
		shelf := NewShelf("Cold shelf", "cold", 5, 1)

		ordersCount := 6
		for i := 0; i < ordersCount; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID}
		}

		ok := shelf.DeleteRandomOrder()

		if !ok {
			t.Errorf("got %v want %v", ok, true)
		}

		if shelf.OrdersCount() != ordersCount-1 {
			t.Errorf("got %v want %v", shelf.OrdersCount(), ordersCount-1)
		}
	})

	t.Run("FindOrderByTemp", func(t *testing.T) {
		shelf := NewShelf("Overflow shelf", "any", 10, 2)

		for i := 0; i < 3; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID, Temperature: "cold"}
		}
		for i := 3; i < 6; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID, Temperature: "hot"}
		}

		order := shelf.FindOrderByTemp("cold")

		if order == nil {
			t.Errorf("got %v want %v", nil, &Order{})
		} else if order.Temperature != "cold" {
			t.Errorf("got %v want %v", order.Temperature, "cold")
		}

		order = shelf.FindOrderByTemp("hot")

		if order == nil {
			t.Errorf("got %v want %v", nil, "Order{}")
		} else if order.Temperature != "hot" {
			t.Errorf("got %v want %v", order.Temperature, "cold")
		}
	})

	t.Run("FindOrderByTemp_Negative", func(t *testing.T) {
		shelf := NewShelf("Overflow shelf", "any", 10, 2)

		for i := 0; i < 3; i++ {
			orderID := fmt.Sprintf("%d", i)
			shelf.orders[orderID] = &Order{ID: orderID, Temperature: "cold"}
		}

		order := shelf.FindOrderByTemp("hot")

		if order != nil {
			t.Errorf("got %v want %v", order, nil)
		}
	})
}

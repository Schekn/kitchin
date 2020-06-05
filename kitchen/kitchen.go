// Package kitchen receives orders and instantly cook the order upon receiving it, and then place the order on the best-available shelf to await pick up by a courier
package kitchen

import (
	c "delivery/config"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Kitchen -
type Kitchen struct {
	// shelves by temperature
	Shelves map[string]*Shelf
	// shelf for orders with any temperature
	OverflowShelf *Shelf

	paused bool
	logger *log.Entry
	mutex  sync.Mutex
}

// New creates new kitchen by given parameters
func New(shelves map[string]*Shelf, overflowShelf *Shelf) *Kitchen {
	capacity := overflowShelf.Capacity
	for _, s := range shelves {
		capacity += s.Capacity
	}

	logger := log.New()

	if os.Getenv("GO_ENV") == "testing" {
		logger.SetOutput(ioutil.Discard)
	}

	k := &Kitchen{
		Shelves:       shelves,
		OverflowShelf: overflowShelf,
		paused:        false,
		logger: logger.WithFields(log.Fields{
			"source":   "kitchen",
			"capacity": capacity,
		}),
	}

	return k
}

// PlaceOrder adds order to the shelf
func (k *Kitchen) PlaceOrder(order *Order) (result bool) {
	shelf, ok := k.Shelves[order.Temperature]
	if !ok {
		k.logger.WithFields(k.getExtraFileds()).Warnf("There is no shelf with temprature '%s' for order with ID %s", order.Temperature, order.ID)
		return false
	}

	result = shelf.AddOrder(order)
	if !result {
		result = k.OverflowShelf.AddOrder(order)
	}

	if !result {
		result = k.RotateOrdersFromOverflowShelve(order)
		if !result {
			k.logger.WithFields(k.getExtraFileds()).Warn("There are no available seats on the kitchen")
		}
	}

	return result
}

// GetAvailableShelves returns available shelves that have empty seats
func (k *Kitchen) GetAvailableShelves() []*Shelf {
	var result []*Shelf

	for _, s := range k.Shelves {
		if s.HasEmptySeats() {
			result = append(result, s)
		}
	}

	return result
}

// IsEmpty checks if the all shelves are empty
func (k *Kitchen) IsEmpty() bool {
	result := true

	for _, s := range k.Shelves {
		if !result {
			break
		}
		result = s.IsEmpty()
	}

	return result
}

// RotateOrdersFromOverflowShelve frees up space on a shelf by moving an order to another shelf or deletes a randomly selected order if there are no free spaces on other shelves
func (k *Kitchen) RotateOrdersFromOverflowShelve(order *Order) (result bool) {
	k.mutex.Lock()

	availableShelves := k.GetAvailableShelves()

	if len(availableShelves) > 0 {
		for _, availableShelf := range availableShelves {
			orderToMove := k.OverflowShelf.FindOrderByTemp(availableShelf.Temperature)
			if orderToMove != nil {
				availableShelf.AddOrder(orderToMove)
				k.OverflowShelf.DeleteOrder(orderToMove.ID)

				result = k.OverflowShelf.AddOrder(order)
				break
			}
		}
	} else {
		k.OverflowShelf.DeleteRandomOrder()
		result = k.OverflowShelf.AddOrder(order)
	}

	k.mutex.Unlock()

	return result
}

// CreateCourier creates a courier for the order
func (k *Kitchen) CreateCourier(order *Order) {
	for {
		rand.Seed(time.Now().UnixNano())
		randValue := rand.Intn(c.Config.Courier.Arrive.Max-c.Config.Courier.Arrive.Min+1) + c.Config.Courier.Arrive.Max
		time.Sleep(time.Duration(randValue) * c.Config.Courier.Arrive.Duration)

		if k.paused {
			continue
		}

		k.logger.WithFields(k.getExtraFileds()).Infof("Courier arrive for order: %s", order.ID)

		ok := k.PickUpOrder(order)
		if ok {
			k.logger.WithFields(k.getExtraFileds()).Infof("Order recived: %s", order.ID)
		} else {
			k.logger.WithFields(k.getExtraFileds()).Warnf("Order not found: %s", order.ID)
		}

		break
	}
}

// PickUpOrder зicks up an order from a shelf
func (k *Kitchen) PickUpOrder(order *Order) (ok bool) {
	_, ok = k.Shelves[order.Temperature].WithdrawOrder(order.ID)
	if !ok {
		_, ok = k.OverflowShelf.WithdrawOrder(order.ID)
	}

	return ok
}

// Pause pauses the kitchen
func (k *Kitchen) Pause() {
	k.paused = true
	k.OverflowShelf.Pause()
	for _, s := range k.Shelves {
		s.Pause()
	}
}

// Unpause unpause the kitchen
func (k *Kitchen) Unpause() {
	k.paused = false
	k.OverflowShelf.Unpause()
	for _, s := range k.Shelves {
		s.Unpause()
	}
}

// IsOnPause returns kitchen state
func (k *Kitchen) IsOnPause() bool {
	return k.paused
}

func (k *Kitchen) getExtraFileds() log.Fields {
	fields := log.Fields{
		"ordersCount": k.OverflowShelf.OrdersCount(),
		"overflowShelf": map[string]int{
			"capacity":    k.OverflowShelf.Capacity,
			"ordersCount": k.OverflowShelf.OrdersCount(),
		},
	}

	for _, s := range k.Shelves {
		fields["ordersCount"] = fields["ordersCount"].(int) + s.OrdersCount()

		fields[s.Temperature+"Shelf"] = map[string]int{
			"capacity":    s.Capacity,
			"ordersCount": s.OrdersCount(),
		}
	}

	return fields
}

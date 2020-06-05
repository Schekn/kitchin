package kitchen

import (
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	c "delivery/config"
)

// Shelf -
type Shelf struct {
	// name of the shelf
	Name string
	// temperature with which orders should be stored on this shelf
	Temperature string
	// max numbers of orders on the shelf
	Capacity int

	decayModifier int
	orders        map[string]*Order
	paused        bool
	logger        *log.Entry
	mutex         sync.Mutex
}

// NewShelf creates new shelve by given parameters
func NewShelf(name, temp string, cap, decayModifier int) *Shelf {
	logger := log.New()

	if os.Getenv("GO_ENV") == "testing" {
		logger.SetOutput(ioutil.Discard)
	}

	shelf := &Shelf{
		Name:        name,
		Temperature: temp,
		Capacity:    cap,

		decayModifier: decayModifier,
		orders:        make(map[string]*Order, cap),
		logger: logger.WithFields(log.Fields{
			"source":      name,
			"temperature": temp,
			"capacity":    cap,
		}),
	}

	go func() {
		for range time.Tick(c.Config.Order.Age.Duration) {
			if shelf.paused {
				continue
			}

			shelf.mutex.Lock()
			for orderID, order := range shelf.orders {
				order.IncAge()
				if order.GetInherentValue() <= 0 {
					delete(shelf.orders, orderID)
					shelf.logger.WithFields(shelf.getExtraFileds()).Warningf("Order expired: %s", orderID)
				}
			}
			shelf.mutex.Unlock()
		}
	}()

	return shelf
}

// OrdersCount returns the count of the orders on the shelf
func (s *Shelf) OrdersCount() int {
	return len(s.orders)
}

// HasEmptySeats checks if the shelf has empty seats
func (s *Shelf) HasEmptySeats() bool {
	return s.Capacity > len(s.orders)
}

// IsEmpty checks if the shelf is empty
func (s *Shelf) IsEmpty() bool {
	return len(s.orders) == 0
}

// AddOrder adds order to the shelf
func (s *Shelf) AddOrder(order *Order) bool {
	s.mutex.Lock()
	if !s.HasEmptySeats() {
		s.mutex.Unlock()
		s.logger.WithFields(s.getExtraFileds()).Warn("There are no empty seats on the shelf")
		return false
	}

	order.shelfDecayModifier = s.decayModifier
	s.orders[order.ID] = order

	s.mutex.Unlock()

	s.logger.Infof("Order added: %s", order.ID)

	return true
}

// WithdrawOrder takes the order off the shelf
func (s *Shelf) WithdrawOrder(orderID string) (order *Order, ok bool) {
	s.mutex.Lock()
	if order, ok = s.orders[orderID]; ok {
		delete(s.orders, orderID)
		s.logger.WithFields(s.getExtraFileds()).Infof("Order withdrawn: %s", order.ID)
	}
	s.mutex.Unlock()

	return order, ok
}

// DeleteOrder deletes the order from the shelf
func (s *Shelf) DeleteOrder(orderID string) bool {
	s.mutex.Lock()
	_, ok := s.orders[orderID]
	if ok {
		delete(s.orders, orderID)
		s.logger.WithFields(s.getExtraFileds()).Infof("Order deleted: %s", orderID)
	}
	s.mutex.Unlock()

	return ok
}

// DeleteRandomOrder deletes the random order from the shelf
func (s *Shelf) DeleteRandomOrder() (result bool) {
	s.mutex.Lock()

	rand.Seed(time.Now().UnixNano())
	randValue := rand.Intn(len(s.orders) - 1)

	i := 0
	for orderID := range s.orders {
		if i == randValue {
			delete(s.orders, orderID)
			s.logger.WithFields(s.getExtraFileds()).Warnf("Order discarded: %s", orderID)
			result = true
			break
		}
		i++
	}

	s.mutex.Unlock()

	return result
}

// FindOrderByTemp returns order from the shelf by given temperature
func (s *Shelf) FindOrderByTemp(temp string) (result *Order) {
	s.mutex.Lock()
	for _, order := range s.orders {
		if order.Temperature == temp {
			result = order
			break
		}
	}
	s.mutex.Unlock()

	return result
}

// Pause pauses the shelf
func (s *Shelf) Pause() {
	s.paused = true
}

// Unpause unpause the shelf
func (s *Shelf) Unpause() {
	s.paused = false
}

func (s *Shelf) getExtraFileds() log.Fields {
	return log.Fields{
		"ordersCount": s.OrdersCount(),
	}
}

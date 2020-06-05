// Real-time system that emulates the fulfillment of delivery orders for a kitchen
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	c "delivery/config"
	"delivery/kitchen"
)

func main() {
	ordersPath := flag.String("o", "", "Orders file path (required)")
	configPath := flag.String("c", "config.yml", "Config file path")
	flag.Parse()

	if *ordersPath == "" {
		flag.PrintDefaults()
		log.Fatal("Missing required arguments")
	}

	err := c.Init(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	orders, err := readOrders(*ordersPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Cannot read orders"))
	}

	log.Infof("%d orders have been read", len(orders))

	k := kitchen.New(
		createShelvesFromConfig(),
		createOverflowShelfFromConfig(),
	)

	command := make(chan string)

	go func() {
		log.Info("Start delivery...")

		for {
			select {
			case cmd := <-command:
				fmt.Println(cmd)
				switch cmd {
				case "p":
					log.Warning("PAUSED")
					k.Pause()
				case "c":
					log.Warning("CONTINUE")
					k.Unpause()
				}
			default:
				if k.IsOnPause() {
					continue
				} else if len(orders) == 0 {
					return
				}

				time.Sleep(c.Config.Order.IngestionRate.Duration / time.Duration(c.Config.Order.IngestionRate.Count))

				var order *kitchen.Order
				order, orders = orders[0], orders[1:]

				log.Infof("Order received: %s", order.ID)

				if k.PlaceOrder(order) {
					go k.CreateCourier(order)
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if k.IsEmpty() && !k.IsOnPause() {
				os.Exit(0)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd != "p" && cmd != "c" {
			if cmd != "" {
				log.Warning("Wrong command")
			}
			continue
		}

		command <- cmd
	}
}

func readOrders(path string) ([]*kitchen.Order, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result []*kitchen.Order

	err = json.Unmarshal([]byte(file), &result)

	return result, err
}

func createShelvesFromConfig() map[string]*kitchen.Shelf {
	shelves := make(map[string]*kitchen.Shelf, len(c.Config.Shelves))
	for _, shelfData := range c.Config.Shelves {
		shelves[shelfData.Temperature] = kitchen.NewShelf(shelfData.Name, shelfData.Temperature, shelfData.Capacity, shelfData.DecayModifier)
	}

	return shelves
}

func createOverflowShelfFromConfig() *kitchen.Shelf {
	return kitchen.NewShelf(
		c.Config.OverflowShelf.Name,
		c.Config.OverflowShelf.Temperature,
		c.Config.OverflowShelf.Capacity,
		c.Config.OverflowShelf.DecayModifier,
	)
}

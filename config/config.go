// Packge config describes configuration for delivery system
package config

import (
	"errors"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// DeliveryConfig -
type DeliveryConfig struct {
	Order struct {
		IngestionRate struct {
			Count    int           `yaml:"count"`
			Time     string        `yaml:"time"`
			Duration time.Duration `yaml:"-"`
		} `yaml:"ingestionRate"`
		Age struct {
			Time     string        `yaml:"time"`
			Duration time.Duration `yaml:"-"`
		} `yaml:"age"`
	} `yaml:"order"`
	Shelves []struct {
		Name          string `yaml:"name"`
		Temperature   string `yaml:"temp"`
		Capacity      int    `yaml:"cap"`
		DecayModifier int    `yaml:"decayModifier​"`
	} `yaml:"shelves"`
	OverflowShelf struct {
		Name          string `yaml:"name"`
		Temperature   string `yaml:"temp"`
		Capacity      int    `yaml:"cap"`
		DecayModifier int    `yaml:"decayModifier​"`
	} `yaml:"overflowShelf"`
	Courier struct {
		Arrive struct {
			Time     string        `yaml:"time"`
			Duration time.Duration `yaml:"-"`
			Min      int           `yaml:"min"`
			Max      int           `yaml:"max"`
		} `yaml:"arrive"`
	} `yaml:"courier"`
}

// Config delivery config
var Config *DeliveryConfig

// Init initialize configuration
func Init(path string) error {
	Config = &DeliveryConfig{}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(contents, Config)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(Config.Order.IngestionRate.Time)
	if err != nil {
		return err
	}

	Config.Order.IngestionRate.Duration = duration

	duration, err = time.ParseDuration(Config.Order.Age.Time)
	if err != nil {
		return err
	}

	Config.Order.Age.Duration = duration

	if Config.Courier.Arrive.Min > Config.Courier.Arrive.Max {
		return errors.New("'Courier.Arrive.Min' cannot be less then 'Courier.Arrive.Max'")
	}

	duration, err = time.ParseDuration(Config.Courier.Arrive.Time)
	if err != nil {
		return err
	}

	Config.Courier.Arrive.Duration = duration

	return nil
}

package kitchen

// Order -
type Order struct {
	// ID of the order
	ID string `json:"id"`
	// name of the order
	Name string `json:"name"`
	// preferred shelf storage temperature
	Temperature string `json:"temp"`
	// shelf wait max duration (seconds)
	ShelfLife int `json:"shelfLife"`
	// value deterioration modifier
	DecayRate float64 `json:"decayRate"`

	shelfDecayModifier int
	age                int
}

// GetInherentValue returns order's have an inherent value that will deteriorate over time, based on the order’s ​shelfLife​ and decayRate​ fields
func (o *Order) GetInherentValue() float64 {
	return (float64(o.ShelfLife) - o.DecayRate*float64(o.age*o.shelfDecayModifier)) / float64(o.ShelfLife)
}

// IncAge increases the age of the order by one
func (o *Order) IncAge() {
	o.age++
}

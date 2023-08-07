package main

// Order an order from a restaurant that serves pizza, beer, and salad
type Order struct {
	Items []*OrderItem
}

type OrderItem struct {
	Pizza      *Pizza       `json:"pizza,omitempty"`
	Beer       *Beer        `json:"beer,omitempty"`
	Salad      *Salad       `json:"salad,omitempty"`
	NamedPizza *NamedPizza  `json:"named_pizza,omitempty"`
	Unknown    *UnknownText `json:"unknown,omitempty"`
}

// UnknownText Use this struct for order items that match nothing else
type UnknownText struct {
	Text string // The text that wasn't understood
}

// SizeEnum Define a custom type SizeEnum with an underlying type of int.
type SizeEnum int

const (
	UnknownSize SizeEnum = 0
	Small       SizeEnum = 1
	Medium      SizeEnum = 2
	Large       SizeEnum = 3
	ExtraLarge  SizeEnum = 4
)

func (s SizeEnum) String() string {
	switch s {
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Large:
		return "large"
	case ExtraLarge:
		return "extra large"
	default:
		return ""
	}
}

// NameEnum Define a custom type NameEnum with an underlying type of int.
type NameEnum int

const (
	UnknownName NameEnum = iota
	Hawaiian
	Yeti
	PigInaForest
	CherryBomb
)

func (n NameEnum) String() string {
	switch n {
	case Hawaiian:
		return "Hawaiian"
	case Yeti:
		return "Yeti"
	case PigInaForest:
		return "Pig In a Forest"
	case CherryBomb:
		return "Cherry Bomb"
	default:
		return ""
	}
}

type Pizza struct {
	Size            SizeEnum `json:"size,omitempty"`             // size use a custom SizeEnum type size with an underlying type of int, default: 3
	AddedToppings   []string `json:"added_toppings,omitempty"`   // toppings requested (examples: pepperoni, arugula)
	RemovedToppings []string `json:"removed_toppings,omitempty"` // toppings requested to be removed (examples: fresh garlic, anchovies)
	Quantity        int      `json:"quantity,omitempty"`         // quantity, default: 1
	Name            NameEnum `json:"name,omitempty"`             // used if the requester references a pizza by name
}

type NamedPizza struct {
	Pizza
}

type Beer struct {
	Kind     string // examples: Mack and Jacks, Sierra Nevada, Pale Ale, Miller Lite
	Quantity int    `json:"quantity,omitempty"` // quantity, default: 1
}

var saladSize = []string{"half", "whole"}

var saladStyle = []string{"Garden", "Greek"}

type Salad struct {
	Portion            string   `json:"portion,omitempty"`             // default: half
	Style              string   `json:"style,omitempty"`               // default: Garden
	AddedIngredients   []string `json:"added_ingredients,omitempty"`   // ingredients requested (examples: parmesan, croutons)
	RemovedIngredients []string `json:"removed_ingredients,omitempty"` // ingredients requested to be removed (example: red onions)
	Quantity           int      `json:"quantity,omitempty"`            // quantity, default: 1
}

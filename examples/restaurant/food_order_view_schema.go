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

// UnknownText Use this type for order items that match nothing else
type UnknownText struct {
	Text string // The text that wasn't understood
}

type Size int

const (
	UnknownSize Size = iota
	Small
	Medium
	Large
	ExtraLarge
)

func (s Size) String() string {
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

type Name int

const (
	UnknownName Name = iota
	Hawaiian
	Yeti
	PigInaForest
	CherryBomb
)

func (n Name) String() string {
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
	Size            string   `json:"size,omitempty"`             // [small | medium | large | extra large], default: large
	AddedToppings   []string `json:"added_toppings,omitempty"`   // toppings requested (examples: pepperoni, arugula)
	RemovedToppings []string `json:"removed_toppings,omitempty"` // toppings requested to be removed (examples: fresh garlic, anchovies)
	Quantity        int      `json:"quantity,omitempty"`         // default: 1
	Name            string   `json:"name,omitempty"`             // used if the requester references a pizza by name [Hawaiian | Yeti | Pig In a Forest | Cherry Bomb]
}

type NamedPizza struct {
	Pizza
}

type Beer struct {
	Kind     string // examples: Mack and Jacks, Sierra Nevada Pale Ale, Miller Lite
	Quantity int    `json:"quantity,omitempty"` // default: 1
}

var saladSize = []string{"half", "whole"}

var saladStyle = []string{"Garden", "Greek"}

type Salad struct {
	Portion            string   `json:"portion,omitempty"`             // default: half
	Style              string   `json:"style,omitempty"`               // default: Garden
	AddedIngredients   []string `json:"added_ingredients,omitempty"`   // ingredients requested (examples: parmesan, croutons)
	RemovedIngredients []string `json:"removed_ingredients,omitempty"` // ingredients requested to be removed (example: red onions)
	Quantity           int      `json:"quantity,omitempty"`            // default: 1
}

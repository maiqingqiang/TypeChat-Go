package main

import (
	"fmt"
	"github.com/maiqingqiang/typechat"
	"log"
	"os"
)

var saladIngredients = []string{
	"lettuce",
	"tomatoes",
	"red onions",
	"olives",
	"peppers",
	"parmesan",
	"croutons",
}

var pizzaToppings = []string{
	"pepperoni",
	"sausage",
	"mushrooms",
	"basil",
	"extra cheese",
	"extra sauce",
	"anchovies",
	"pineapple",
	"olives",
	"arugula",
	"Canadian bacon",
	"Mama Lil's Peppers",
}

var namedPizzas = map[string][]string{
	"Hawaiian":        {"pineapple", "Canadian bacon"},
	"Yeti":            {"extra cheese", "extra sauce"},
	"Pig In a Forest": {"mushrooms", "basil", "Canadian bacon", "arugula"},
	"Cherry Bomb":     {"pepperoni", "sausage", "Mama Lil's Peppers"},
}

func main() {
	model, err := typechat.NewLanguageModel()
	if err != nil {
		log.Fatal(err)
	}

	schema, err := os.ReadFile("food_order_view_schema.go")
	if err != nil {
		log.Fatalf("os.ReadFile Error: %v\n", err)
	}

	translator := typechat.NewJsonTranslator[Order](model, string(schema), "Order")

	_ = typechat.ProcessRequests("🍕> ", os.Args[1], func(request string) error {
		order, err := translator.Translate(request)
		if err != nil {
			log.Fatalf("translator.Translate Error: %v\n", err)
		}

		printOrder(order)

		return nil
	})
}

func printOrder(order *Order) {
	if order != nil && len(order.Items) > 0 {
		for _, item := range order.Items {
			if item.Unknown != nil {
				break
			}

			if item.Pizza != nil || item.NamedPizza != nil {
				if item.Pizza.Name != "" {
					addedToppings, ok := namedPizzas[item.Pizza.Name]
					if ok {
						if item.Pizza.AddedToppings != nil {
							item.Pizza.AddedToppings = append(item.Pizza.AddedToppings, addedToppings...)
						} else {
							item.Pizza.AddedToppings = addedToppings
						}
					}
				}

				if item.Pizza.Size == "" {
					item.Pizza.Size = "large"
				}

				quantity := 1
				if item.Pizza.Quantity > 0 {
					quantity = item.Pizza.Quantity
				}

				pizzaStr := fmt.Sprintf(`    %d %s pizza`, quantity, item.Pizza.Size)

				if len(item.Pizza.AddedToppings) > 0 && len(item.Pizza.RemovedToppings) > 0 {
					item.Pizza.AddedToppings, item.Pizza.RemovedToppings = removeCommonStrings(item.Pizza.AddedToppings, item.Pizza.RemovedToppings)
				}

				if len(item.Pizza.AddedToppings) > 0 {
					pizzaStr += " with"
					for index, addedTopping := range item.Pizza.AddedToppings {
						if contains(pizzaToppings, addedTopping) {
							if index == 0 {
								pizzaStr += fmt.Sprintf(" %s", addedTopping)
							} else {
								pizzaStr += fmt.Sprintf(", %s", addedTopping)
							}
						} else {
							log.Printf("We are out of %s", addedTopping)
						}
					}
				}

				if len(item.Pizza.RemovedToppings) > 0 {
					pizzaStr += " and without"
					for index, removedTopping := range item.Pizza.RemovedToppings {
						if index == 0 {
							pizzaStr += fmt.Sprintf(" %s", removedTopping)
						} else {
							pizzaStr += fmt.Sprintf(", %s", removedTopping)
						}
					}
				}
			} else if item.Beer != nil {
				quantity := 1
				if item.Beer.Quantity > 0 {
					quantity = item.Beer.Quantity
				}

				beerStr := fmt.Sprintf("    %d %s", quantity, item.Beer.Kind)
				log.Printf(beerStr)
			} else if item.Salad != nil {
				quantity := 1
				if item.Salad.Quantity > 0 {
					quantity = item.Salad.Quantity
				}

				if item.Salad.Portion != "" {
					item.Salad.Portion = "half"
				}

				if item.Salad.Style != "" {
					item.Salad.Style = "Garden"
				}

				saladStr := fmt.Sprintf(`    %d %s %s salad`, quantity, item.Salad.Portion, item.Salad.Style)

				if len(item.Salad.AddedIngredients) > 0 && len(item.Salad.RemovedIngredients) > 0 {
					item.Salad.AddedIngredients, item.Salad.RemovedIngredients = removeCommonStrings(item.Salad.AddedIngredients, item.Salad.RemovedIngredients)
				}

				if len(item.Salad.AddedIngredients) > 0 {
					saladStr += " with"
					for index, addedIngredient := range item.Salad.AddedIngredients {
						if contains(saladIngredients, addedIngredient) {
							if index == 0 {
								saladStr += fmt.Sprintf(" %s", addedIngredient)
							} else {
								saladStr += fmt.Sprintf(", %s", addedIngredient)
							}
						} else {
							log.Printf("We are out of %s", addedIngredient)
						}
					}
				}

				if len(item.Salad.RemovedIngredients) > 0 {
					saladStr += " and without"
					for index, removedIngredient := range item.Salad.RemovedIngredients {
						if index == 0 {
							saladStr += fmt.Sprintf(" %s", removedIngredient)
						} else {
							saladStr += fmt.Sprintf(", %s", removedIngredient)
						}
					}
				}

				log.Printf(saladStr)
			}
		}
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func removeCommonStrings(a, b []string) ([]string, []string) {
	aSet := make(map[string]struct{})
	for _, item := range a {
		aSet[item] = struct{}{}
	}

	bSet := make(map[string]struct{})
	for _, item := range b {
		bSet[item] = struct{}{}
	}

	for item := range aSet {
		if _, ok := bSet[item]; ok {
			delete(aSet, item)
			delete(bSet, item)
		}
	}

	var aResult []string
	for item := range aSet {
		aResult = append(aResult, item)
	}

	var bResult []string
	for item := range bSet {
		bResult = append(bResult, item)
	}

	return aResult, bResult
}

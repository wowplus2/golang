package main

import (
	"fmt"
	"strings"
)

type Items []Coster

type Itemer interface {
	Coster
	fmt.Stringer
}

type Coster interface {
	Cost() float64
}

type Item struct {
	name  string
	price float64
	qty   int
}

type RentalPeriod int

type DcItem struct {
	Item
	dcRate float64
}

type Rental struct {
	name       string
	freePerDay float64
	periodLen  int
	RentalPeriod
}

type Order struct {
	Itemer
	taxRate float64
}

const (
	Days RentalPeriod = iota
	Weeks
	Months
)

func (p RentalPeriod) ToDays() int {
	switch p {
	case Weeks:
		return 7
	case Months:
		return 30
	default:
		return 1
	}
}

func (t Item) Cost() float64 {
	return t.price * float64(t.qty)
}

func (t DcItem) Cost() float64 {
	return t.Item.Cost() * (1.0 - t.dcRate/100)
}

func (r Rental) Cost() float64 {
	return r.freePerDay * float64(r.ToDays()*r.periodLen)
}

func (ts Items) Cost() (c float64) {
	for _, t := range ts {
		c += t.Cost()
	}
	return
}

func (o Order) Cost() float64 {
	return o.Itemer.Cost() * (1.0 + o.taxRate/100)
}

func display(c Coster) {
	fmt.Printf("cost: %.0f", c.Cost())
}

func (t Item) String() string {
	return fmt.Sprintf("[%s] %.0f", t.name, t.Cost())
}

func (t DcItem) String() string {
	return fmt.Sprintf("%s => %.0f(%.0f%s DC)", t.Item.String(), t.Cost(), t.dcRate, "%")
}

func (r Rental) String() string {
	return fmt.Sprintf("[%s] %.0f", r.name, r.Cost())
}

func (ts Items) String() string {
	var s []string
	for _, v := range ts {
		s = append(s, fmt.Sprint(v))
	}

	return fmt.Sprintf("%d items. total: %.0f\n\t- %s", len(s), ts.Cost(), strings.Join(s, "\n\t- "))
}

func (o Order) String() string {
	return fmt.Sprintf("\nTotal price: %0.f(tax rate: %.2f)\n\tOrder details: %s", o.Cost(), o.taxRate, o.Itemer.String())
}

func main() {
	shirt := Item{"Men's Slim-Fit Shirt", 25000, 3}
	video := Rental{"Interstellar", 1000, 3, Days}
	eShoes := DcItem{
		Item{"Women's Walking Shoes", 50000, 3},
		10.00,
	}

	order1 := Order{Items{shirt, eShoes}, 10.00}
	order2 := Order{video, 5.00}

	fmt.Println(order1)
	fmt.Println(order2)
}

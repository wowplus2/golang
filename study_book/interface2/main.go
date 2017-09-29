package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// 제너릭 컬렉션
type Items []Coster

type Coster interface {
	Cost() float64
}

type Item struct {
	name     string
	price    float64
	quantity int
}

type DiscountItem struct {
	Item
	dcRate float64
}

type Rental struct {
	name       string
	freePerDay float64
	periodLen  int
	RentalPeriod
}

type RentalPeriod int

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

func (i Item) Cost() float64 {
	return i.price * float64(i.quantity)
}

func (t DiscountItem) Cost() float64 {
	return t.Item.Cost() * (1.0 - t.dcRate/100)
}

func (r Rental) Cost() float64 {
	return r.freePerDay * float64(r.ToDays()*r.periodLen)
}

// 제너릭 컬렉션
func (ts Items) Cost() (c float64) {
	for _, t := range ts {
		c += t.Cost()
	}
	return
}

func dispCost(c Coster) {
	fmt.Printf("cost: %f\n", c.Cost())
}

func (t Item) String() string {
	return fmt.Sprintf("[%s] %.0f", t.name, t.Cost())
}

func (t DiscountItem) String() string {
	return fmt.Sprintf("%s => %.0f(%.0f%s DC)", t.Item.String(), t.Cost(), t.dcRate, "%")
}

func (t Rental) String() string {
	return fmt.Sprintf("[%s] %.0f", t.name, t.Cost())
}

func (ts Items) String() string {
	var s []string
	for _, t := range ts {
		s = append(s, fmt.Sprint(t))
	}
	return fmt.Sprintf("%d items. total: %0.f\n\t -%s", len(s), ts.Cost(), strings.Join(s, "\n\t- "))
}

func haldle(w io.Writer, msg string) {
	fmt.Fprint(w, msg)
}

func main() {
	//shoes := Item{"Sport Shoes", 30000, 2}
	dcShoes := DiscountItem{
		Item{"Women's Walking Shoes", 50000, 3},
		10.00,
	}
	shirt := Item{"Men's Slim-Fit Shirt", 25000, 3}
	video := Rental{"Installer", 1000, 3, Days}

	// 제너릭 컬렉션
	items := Items{shirt, video, dcShoes}

	fmt.Println(shirt)
	fmt.Println(video)
	fmt.Println(dcShoes)
	fmt.Println(items)
	//dispCost(items)
	/*
		dispCost(shoes)
		dispCost(dcShoes)
		dispCost(shirt)
		dispCost(video)
	*/

	msg := []string{"This", "is", "an", "example", "of", "io.Writer.", "...", "...", "...", "...", "...", "...", "..."}

	for _, s := range msg {
		time.Sleep(100 * time.Millisecond)
		haldle(os.Stdout, s+" ")
	}
}

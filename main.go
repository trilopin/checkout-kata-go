package main

import (
	"fmt"
	"os"
)

// Item is the representation of a product in the inventory
type Item struct {
	Code  string
	Name  string
	Price float64
}

// DiscountApplier is a common interface for discount implementations
type DiscountApplier interface {
	// ApplyDiscount computes total discount from current units adn price
	ApplyDiscount(n int, price float64) float64
}

// GetFreeUnitDiscount is a "Buy X -> get Y Free units" dicount
type GetFreeUnitDiscount struct {
	Buy  int
	Free int
}

// ApplyDiscount returns free units and same price
func (gfud *GetFreeUnitDiscount) ApplyDiscount(n int, price float64) float64 {
	free := int(n / gfud.Buy)
	discount := float64(free) * price
	return discount
}

// LessPriceDiscount is a "fixed lower prices after buying X units" discount
type LessPriceDiscount struct {
	Buy   int
	Price float64
}

// ApplyDiscount computes total discount and returns 0 free items
func (lpd *LessPriceDiscount) ApplyDiscount(n int, price float64) float64 {
	discount := 0.0
	if n >= lpd.Buy {
		discount = (price - lpd.Price) * float64(n)
	}
	return discount
}

// PricingRules is the complete inventory plus special discounts
type PricingRules struct {
	products  map[string]Item
	discounts map[string]DiscountApplier
}

// NewPricingRules builds and initialises a PricingRules struct
func NewPricingRules() *PricingRules {
	pr := &PricingRules{}
	pr.products = make(map[string]Item)
	pr.discounts = make(map[string]DiscountApplier)
	return pr
}

// Add makes a provision for one or more products
func (pr *PricingRules) Add(items ...Item) error {
	for _, item := range items {
		_, ok := pr.products[item.Code]
		if !ok {
			pr.products[item.Code] = item
		} else {
			return fmt.Errorf("item %s already exists", item.Code)
		}
	}
	return nil
}

// RegisterDiscount setups a discount for an item
func (pr *PricingRules) RegisterDiscount(code string, d DiscountApplier) error {
	if _, ok := pr.products[code]; !ok {
		return fmt.Errorf("item %s does not exists", code)
	}
	if _, ok := pr.discounts[code]; ok {
		return fmt.Errorf("discount for item %s already exists", code)
	}
	pr.discounts[code] = d
	return nil
}

// Checkout holds rules and products
type Checkout struct {
	pricingRules *PricingRules
	products     map[string]int
}

// NewCheckout creates and initialises Checkout struct
func NewCheckout(pr *PricingRules) *Checkout {
	c := &Checkout{pricingRules: pr}
	c.products = make(map[string]int)
	return c
}

func (c *Checkout) getPrice(code string) float64 {
	product, _ := c.pricingRules.products[code]
	return product.Price
}

// Scan add a product to the basket
func (c *Checkout) Scan(code string) error {
	_, ok := c.pricingRules.products[code]
	if !ok {
		return fmt.Errorf("item %s does not exists in inventory", code)
	}
	c.products[code]++
	return nil
}

// GetTotal computes final price for all scanned products
func (c *Checkout) GetTotal() float64 {
	total := 0.0
	for code, amount := range c.products {
		price := c.getPrice(code)
		total += float64(amount) * price
		discount, ok := c.pricingRules.discounts[code]
		if ok {
			amountDiscount := discount.ApplyDiscount(amount, price)
			total -= amountDiscount
		}
	}
	return total
}

func main() {
	var items = map[string]Item{
		"VOUCHER": Item{
			Code:  "VOUCHER",
			Name:  "Voucher",
			Price: 5.00,
		},
		"TSHIRT": Item{
			Code:  "TSHIRT",
			Name:  "T-Shirt",
			Price: 20.00,
		},
		"MUG": Item{
			Code:  "MUG",
			Name:  "Coffee Mug",
			Price: 7.50,
		},
	}
	pr := NewPricingRules()
	for _, item := range items {
		pr.Add(item)
	}
	pr.RegisterDiscount("VOUCHER", &GetFreeUnitDiscount{2, 1})
	pr.RegisterDiscount("TSHIRT", &LessPriceDiscount{3, 19.00})

	co := NewCheckout(pr)
	for _, product := range os.Args[1:] {
		if err := co.Scan(product); err != nil {
			fmt.Println("Unknown product " + product)
			return
		}
	}
	total := co.GetTotal()
	fmt.Printf("\nTotal Price: %0.2f€\n", total)
}

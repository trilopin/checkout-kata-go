package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestNewPricingRules(t *testing.T) {
	pr := NewPricingRules()
	assert.NotNil(t, pr)
	assert.Equal(t, len(pr.products), 0)
}

func TestPricingRuleAdd(t *testing.T) {
	pr := NewPricingRules()
	err := pr.Add(items["VOUCHER"])
	assert.Equal(t, len(pr.products), 1)
	assert.Nil(t, err)
}
func TestPricingRuleAddFails(t *testing.T) {
	pr := NewPricingRules()
	err := pr.Add(items["VOUCHER"], items["VOUCHER"])
	assert.Equal(t, len(pr.products), 1)
	assert.Equal(t, err, errors.New("item VOUCHER already exists"))
}

// createPricingRules is a helper function for checkout testcases
func createPricingRules() *PricingRules {
	pr := NewPricingRules()
	for _, item := range items {
		pr.Add(item)
	}
	pr.RegisterDiscount("VOUCHER", &GetFreeUnitDiscount{2, 1})
	pr.RegisterDiscount("TSHIRT", &LessPriceDiscount{3, 19.00})
	return pr
}

func TestNewCheckout(t *testing.T) {
	pricingRules := &PricingRules{}
	checkout := NewCheckout(pricingRules)
	assert.NotNil(t, checkout)
	assert.Equal(t, checkout.pricingRules, pricingRules)
}

func TestCheckoutScan(t *testing.T) {

	var err error
	var tests = []struct {
		items    []string
		expected map[string]int
		err      error
	}{
		{
			[]string{"VOUCHER"},
			map[string]int{"VOUCHER": 1},
			nil,
		},
		{
			[]string{"VOUCHER", "VOUCHER"},
			map[string]int{"VOUCHER": 2},
			nil,
		},
		{
			[]string{"VOUCHER", "FAKE"},
			map[string]int{"VOUCHER": 1},
			errors.New("item FAKE does not exists in inventory"),
		},
	}

	for _, test := range tests {
		assert.Nil(t, err)
		checkout := NewCheckout(createPricingRules())
		for _, item := range test.items {
			err = checkout.Scan(item)
		}
		assert.Equal(t, err, test.err)
		assert.Equal(t, checkout.products, test.expected)
	}
}

func TestCheckoutSGetTotal(t *testing.T) {
	var err error
	var tests = []struct {
		items    []string
		expected float64
	}{
		{
			[]string{},
			0.00,
		},
		{
			[]string{"VOUCHER"},
			5.00,
		},
		{
			[]string{"VOUCHER", "TSHIRT", "VOUCHER"},
			25.00,
		},
		{
			[]string{"VOUCHER", "TSHIRT", "MUG"},
			32.50,
		},
		{
			[]string{"TSHIRT", "TSHIRT", "TSHIRT", "VOUCHER", "TSHIRT"},
			81.00,
		},
		{
			[]string{"VOUCHER", "TSHIRT", "VOUCHER", "VOUCHER", "MUG", "TSHIRT", "TSHIRT"},
			74.50,
		},
	}

	for _, test := range tests {
		assert.Nil(t, err)
		checkout := NewCheckout(createPricingRules())
		for _, item := range test.items {
			checkout.Scan(item)
		}
		total := checkout.GetTotal()
		assert.Equal(t, total, test.expected)
	}
}

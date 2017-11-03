# checkout-kata

## Problem
A physical store which sells 3 products:

```
Code         | Name                |  Price
-------------------------------------------------
VOUCHER      | Voucher      |   5.00€
TSHIRT       | T-Shirt      |  20.00€
MUG          | Cafify Coffee Mug   |   7.50€
```

Various departments have insisted on the following discounts:

 * The marketing department thinks a buy 2 get 1 free promotions will work best (buy two of the same product, get an additional one free), and would like this to only apply to `VOUCHER` items.

 * The CFO insists that the best way to increase sales is with discounts on bulk purchases (buying x or more of a product, the price of that product is reduced), and requests that if you buy 3 or more `TSHIRT` items, the price per unit should be 19.00€.

 * Checkout process allows for items to be scanned in any order, and should return the total amount to be paid. The interface for the checkout process look like this (GO):

```golang
co := NewCheckout(pricingRules)
co.Scan("VOUCHER")
co.Scan("VOUCHER")
co.Scan("TSHIRT")
price := co.GetTotal()
```

Examples:

    Items: VOUCHER, TSHIRT, MUG
    Total: 32.50€

    Items: VOUCHER, TSHIRT, VOUCHER
    Total: 25.00€

    Items: TSHIRT, TSHIRT, TSHIRT, VOUCHER, TSHIRT
    Total: 81.00€

    Items: VOUCHER, TSHIRT, VOUCHER, VOUCHER, MUG, TSHIRT, TSHIRT
    Total: 74.50€

## Solution

### Install
Install testify/assert package
```
$ dep ensure
```

Run tests
```
$ go test --cover
PASS
coverage: 72.7% of statements
ok      github.com/trilopin/checkout-kata       0.010s
```

Execute
```
$ go run main.go VOUCHER VOUCHER TSHIRT MUG
Total Price: 32.50€
```

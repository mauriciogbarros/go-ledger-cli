package currency

import "fmt"

type Currency int64

func (c Currency) String() string {
	return fmt.Sprintf("%.2f", float64(c)/100)
}

func Convert64(value float64) Currency {
	var cents Currency = Currency(value * 100)
	if cents % 2 != 0 {
		cents += 1
	}
	return cents
}
package types

import "fmt"

type Currency int64

func (c Currency) String() string {
	return fmt.Sprintf("%.2f", float64(c)/100)
}

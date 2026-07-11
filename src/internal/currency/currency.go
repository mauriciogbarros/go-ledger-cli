package currency

import (
	"fmt"
	"strings"
)

type Currency int64

func (c Currency) String() string {
	var output string
	var num string = fmt.Sprintf("%d.%02d", c / 100, c % 100)
	var parts []string = strings.Split(num, ".")
	var fPart, iPart string
		if len(parts) > 1 {
		fPart = "." + parts[1]
	} else {
		fPart = ".00"
	}
	iPart = parts[0]
	for len(iPart) > 0 {
		if len(iPart) > 3 {
			output = "," + iPart[len(iPart)-3:] + output
			iPart = iPart[:len(iPart)-3]
		} else {
			output = iPart + output
			iPart = ""
		}
	}
	return output + fPart
}

func Convert64(value float64) Currency {
	var thous int64 = int64(value * 1000)
	var tcents int64 = int64(value * 100) * 10
	var dif = thous - tcents
	var cents int64 = 0
	if dif < 5 {
		cents = tcents / 10
	}
	if dif > 5 {
		cents = (tcents + 10) / 10
	}
	if dif == 5 {
		if (tcents / 10) % 2 == 0 {
			cents = tcents / 10
		} else {
			cents = (tcents + 10) / 10
		}
	}
	return Currency(cents)
}
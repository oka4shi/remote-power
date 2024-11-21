package main

import (
	"github.com/oka4shi/remote-power/gpio"
)

func main() {
	p, err := gpio.NewPort(gpio.BANK_3, gpio.GROUP_C, gpio.X_5, gpio.OUT)
	if err != nil {
		panic("やばい!")
	}

	p.DegitalWrite(gpio.HIGH)
}

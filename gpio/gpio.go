package gpio

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Port struct {
	Pin int
}

func NewPort(bank Bank, group Group, x X) (*Port, error) {
	pin := int(bank)*32 + (int(group)*8 + int(x))
	p := Port{
		Pin: pin,
	}
	export, err := os.Stat("/sys/class/gpio/export")
	if err != nil || export.IsDir() {
		return nil, errors.New("Can't use GPIO")
	}

	_, err = os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d", pin))
	if err == nil {
		return &p, nil
	}

	if err := os.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(pin)), 0644); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to initialize gpio%d", pin))
	}
	return &p, nil
}

func (p *Port) SetDirection(value Direction) error {

	// TODO
	return nil
}

func (p *Port) DegitalWrite(value OutLevel) error {

	// TODO
	return nil
}

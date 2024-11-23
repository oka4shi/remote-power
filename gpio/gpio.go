package gpio

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Port struct {
	Mu     sync.Mutex
	Pin    int
	Locked bool
}

func isPortReady(pin int) bool {
	gpioFile, err := os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d", pin))
	return err == nil && gpioFile.IsDir()
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

	if isPortReady(pin) {
		return &p, nil
	}

	if err := os.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(pin)), 0644); err != nil {
		return nil, errors.New("Failed to initialize GPIO")
	}
	return &p, nil
}

var directionMap = map[Direction]string{
	IN:  "in",
	OUT: "out",
}

func (p *Port) Lock() error {
	p.Mu.Lock()
	if p.Locked {
		return errors.New("Already locked")
	}

	p.Locked = true
	p.Mu.Unlock()
	return nil
}

func (p *Port) Unlock() {
	p.Mu.Lock()
	p.Locked = false
	p.Mu.Unlock()
}

func (p *Port) SetDirection(value Direction) error {
	if !isPortReady(p.Pin) {
		return errors.New("The port isn't initialized")
	}
	if err := os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Pin), []byte(directionMap[value]), 0644); err != nil {
		return errors.New("Failed to set the direction")
	}

	return nil
}

func (p *Port) DegitalWrite(value OutLevel) error {
	if !isPortReady(p.Pin) {
		return errors.New("The port isn't initialized")
	}
	if err := os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Pin), []byte(fmt.Sprintf("%d", value)), 0644); err != nil {
		return errors.New("Failed to set value")
	}

	return nil
}

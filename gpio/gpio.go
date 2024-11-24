package gpio

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Port struct {
	Mu     sync.Mutex
	Pin    int
	Locked bool
}

var DryRun = false

func writeSysFile(path string, data string) error {
	if DryRun {
		log.Printf("Write \"%s\" to %s", data, path)
		return nil
	}

	err := os.WriteFile(path, []byte(data), 0644)
	return err
}

func doesFileExist(path string) bool {
	if DryRun {
		log.Printf("Check existance of %s", path)
		return true
	}

	gpioFile, err := os.Stat(path)
	return err == nil && gpioFile.IsDir()

}

func isPortReady(pin int) bool {
	return doesFileExist(fmt.Sprintf("/sys/class/gpio/gpio%d", pin))
}

func NewPort(bank Bank, group Group, x X) (*Port, error) {
	pin := int(bank)*32 + (int(group)*8 + int(x))
	p := Port{
		Pin: pin,
	}
	if DryRun {
		return &p, nil
	}

	if !doesFileExist("/sys/class/gpio/export") {
		return nil, errors.New("Can't use GPIO")
	}

	if isPortReady(pin) {
		return &p, nil
	}

	if err := writeSysFile("/sys/class/gpio/export", strconv.Itoa(pin)); err != nil {
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
	defer p.Mu.Unlock()
	if p.Locked {
		return errors.New("Already locked")
	}

	p.Locked = true
	return nil
}

func (p *Port) Unlock() {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	p.Locked = false
}

func (p *Port) SetDirection(value Direction) error {
	if !isPortReady(p.Pin) {
		return errors.New("The port isn't initialized")
	}
	if err := writeSysFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Pin), directionMap[value]); err != nil {
		return errors.New("Failed to set the direction")
	}

	return nil
}

func (p *Port) DegitalWrite(value OutLevel) error {
	if !isPortReady(p.Pin) {
		return errors.New("The port isn't initialized")
	}
	if err := writeSysFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Pin), fmt.Sprintf("%d", value)); err != nil {
		return errors.New("Failed to set value")
	}

	return nil
}

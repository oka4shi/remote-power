package network

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Iface struct {
	Name   string
	Status NetworkStatus
}

type Ifaces = []*Iface

func NewIface(name string) (*Iface, error) {
	status, err := getIfaceStatus(name)
	if err != nil {
		return nil, err
	}

	port := Iface{
		Name:   name,
		Status: status,
	}
	return &port, nil

}

func NewIfaces() (Ifaces, error) {
	ifaces, err := ListPhisicalIfaces("/sys/devices")
	if err != nil {
		return nil, err
	}

	ports := Ifaces{}
	for _, iface := range ifaces {
		iface, err := NewIface(iface)
		if err != nil {
			return nil, err
		}
		ports = append(ports, iface)
	}

	return ports, nil

}

func ListPhisicalIfaces(root string) ([]string, error) {
	findList := []string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "net" || strings.Contains(path, "virtual") {
			return nil
		}

		devices, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, dev := range devices {
			findList = append(findList, dev.Name())
		}

		return nil

	})
	slices.Sort(findList)
	return findList, err
}

func getIfaceStatus(name string) (NetworkStatus, error) {
	file, err := os.Open(filepath.Join("/sys/class/net", name, "carrier"))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	switch string(content) {
	case "0":
		return Down, nil
	case "1":
		return Up, nil
	default:
		return 0, errors.New(fmt.Sprintf("Couldn't get the status of network interface \"%s\"", name))
	}
}

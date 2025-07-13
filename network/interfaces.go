package network

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Iface struct {
	Name   string
	Status NetworkStatus
}

type Ifaces []*Iface

func GetIface(name string) *Iface {
	status, err := getIfaceStatus(name)
	if err != nil {
		log.Println("Error getting status for interface", name, ":", err)
	}

	port := Iface{
		Name:   name,
		Status: status,
	}
	return &port

}

func GetIfaces() (Ifaces, error) {
	ifaces, err := ListPhisicalIfaces("/sys/devices")
	if err != nil {
		return nil, err
	}

	ports := Ifaces{}
	for _, iface := range ifaces {
		iface := GetIface(iface)
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
		return Unknown, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return Unknown, err
	}

	if len(content) == 0 {
		return Unknown, fmt.Errorf("Couldn't get the status of network interface \"%s\"", name)
	}

	status := strings.TrimSpace(string(content))
	switch status {
	case "0":
		return Down, nil
	case "1":
		return Up, nil
	default:
		return Unknown, fmt.Errorf("Couldn't get the status of network interface \"%s\"", name)
	}
}

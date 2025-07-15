//go:build ignore

package main

import (
	"fmt"
	"github.com/oka4shi/remote-power/network"
)

func main() {
	printIfaces()

	c := make (chan []byte)
	go network.WatchNetlinkSocket(c)

	for {
		_ = <- c
		printIfaces()
	}
}

func printIfaces () {
	// Create a new netlink client
	interfaces, err := network.GetIfaces()
	if err != nil {
		panic(err)
	}

	// Print the interfaces
	for _, iface := range interfaces {
		fmt.Printf("%s: %d\n", iface.Name, iface.Status)
	}
}

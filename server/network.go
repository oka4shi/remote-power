package server

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/oka4shi/remote-power/network"
)

type Ifaces struct {
	mu         sync.Mutex      `json:"-"`
	Ifaces     *network.Ifaces `json:"ifaces"`
	UpdateTime time.Time       `json:"update_time"`
	Error      bool            `json:"error"`
}

var ifaces = &Ifaces{}

type Channels struct {
	mu       sync.Mutex
	Channels map[uint64]chan []byte
}

var clients = &Channels{
	mu:       sync.Mutex{},
	Channels: make(map[uint64]chan []byte),
}

func WatchIfaces() {
	interfaces, err := network.GetIfaces()
	if err != nil {
		panic(err)
	}
	updateIfaces(interfaces, err)

	c := make(chan []byte)
	go network.WatchNetlinkSocket(c)

	for {
		_ = <-c

		updateIfaces(network.GetIfaces())
		body, err := json.Marshal(ifaces)
		if err != nil {
			fmt.Println("Error marshalling interfaces:", err)
			continue
		}

		for {
			clients.mu.Lock()
			if len(clients.Channels) == 0 {
				clients.mu.Unlock()
				break
			}

			for id, ch := range clients.Channels {
				select {
				case ch <- body:
				default:
				}
				close(clients.Channels[id])
				delete(clients.Channels, id)
				break // Only exccute for the first channel
			}
			clients.mu.Unlock()
		}
	}
}

func updateIfaces(interfaces network.Ifaces, err error) {
	ifaces.mu.Lock()
	defer ifaces.mu.Unlock()
	if err != nil {
		fmt.Println("Error updating interfaces:", err)

		emptyIfaces := network.Ifaces{}
		ifaces.Ifaces = &emptyIfaces

		ifaces.UpdateTime = time.Now()
		ifaces.Error = true
	} else {
		ifaces.Ifaces = &interfaces
		ifaces.UpdateTime = time.Now()
		ifaces.Error = false
	}
}

func NetworkStatus(w http.ResponseWriter, r *http.Request) {
	ifaces.mu.Lock()
	body, err := json.Marshal(ifaces)
	ifaces.mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("An error occured:%s", err))
	}

	w.Header().Set("Content-Type", "application/json")
	io.Writer.Write(w, body)
}

func NetworkWatch(w http.ResponseWriter, r *http.Request) {
	clientId := rand.Uint64N(1 << 63)
	ch := make(chan []byte)

	clients.mu.Lock()
	clients.Channels[clientId] = ch
	clients.mu.Unlock()

	select {
	case content := <-ch:
		w.Header().Set("Content-Type", "application/json")
		io.Writer.Write(w, content)
	case <-r.Context().Done():
		clients.mu.Lock()
		delete(clients.Channels, clientId)
		clients.mu.Unlock()
	}
}

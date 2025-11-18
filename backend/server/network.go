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

type channelType int

const (
	ChannelTypeOnce channelType = iota
	ChannelTypeStream
)

type Ifaces struct {
	mu         sync.Mutex      `json:"-"`
	Ifaces     *network.Ifaces `json:"ifaces"`
	UpdateTime time.Time       `json:"updated_at"`
	Error      bool            `json:"error"`
}

var ifaces = &Ifaces{}

type channels struct {
	mu       sync.Mutex
	Channels map[uint64]*channel
}

type channel struct {
	Ch   chan []byte
	Type channelType
}

var clients = &channels{
	mu:       sync.Mutex{},
	Channels: make(map[uint64]*channel),
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

		// Broadcast the message to every client
		clients.mu.Lock()
		for id, ch := range clients.Channels {
			if ch == nil {
				delete(clients.Channels, id)
				continue
			}

			select {
			case ch.Ch <- body:
			default:
			}

			if ch.Type == ChannelTypeOnce {
				close(ch.Ch)
				delete(clients.Channels, id)
			}
		}
		clients.mu.Unlock()
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
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, fmt.Sprintf("An error occured:%s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Writer.Write(w, body)
}

func NetworkWatch(w http.ResponseWriter, r *http.Request) {
	clientId := rand.Uint64N(1 << 63)
	ch := make(chan []byte)

	clients.mu.Lock()
	clients.Channels[clientId] = &channel{
		Ch:   ch,
		Type: ChannelTypeOnce,
	}
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

func NetworkWatchStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "An error has occured")
		return
	}

	clientId := rand.Uint64N(1 << 63)
	ch := make(chan []byte)

	clients.mu.Lock()
	clients.Channels[clientId] = &channel{
		Ch:   ch,
		Type: ChannelTypeStream,
	}
	clients.mu.Unlock()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	pingTicker := time.NewTicker(1 * time.Second)

	for {
		select {
		case content := <-ch:
			io.WriteString(w, "data: ")
			io.Writer.Write(w, content)
			io.WriteString(w, "\n\n")
			flusher.Flush()
		case <-pingTicker.C:
			fmt.Fprintf(w, "event: heartbeat\n")
			ifaces.mu.Lock()
			fmt.Fprintf(w, "data: %s\n\n", ifaces.UpdateTime.Format(time.RFC3339))
			ifaces.mu.Unlock()
			flusher.Flush()
		case <-r.Context().Done():
			clients.mu.Lock()
			// Not close channel from receiver side to avoid panic
			// Just delete it from the map and let the GC do the rest
			delete(clients.Channels, clientId)
			clients.mu.Unlock()
			return
		}
	}
}

package server

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"

	"github.com/oka4shi/remote-power/gpio"
)

const tokenLength = 8

type Process struct {
	wg    *sync.WaitGroup
	token string
	err   error
}

var process = &Process{}

var p *gpio.Port

func init() {
	port, err := gpio.NewPort(gpio.BANK_3, gpio.GROUP_C, gpio.X_5)
	if err != nil {
		panic("やばい!")
	}
	if err := p.SetDirection(gpio.OUT); err != nil {
		panic("やばい！")
	}

	p = port
}

func pushButton(process *Process, isLong bool) {
	defer (func() {
		process.wg.Done()
		p.Unlock()
	})()
	var duration time.Duration
	if isLong {
		duration = 5 * time.Second
	} else {
		duration = 200 * time.Millisecond
	}

	if err := p.DegitalWrite(gpio.HIGH); err != nil {
		process.err = err
		return
	}

	time.Sleep(duration)

	if err := p.DegitalWrite(gpio.LOW); err != nil {
		process.err = err
		return
	}
}

func Push(w http.ResponseWriter, r *http.Request) {
	isLong := len(r.URL.Query().Get("long")) != 0

	// Lock a GPIO port(return 409 error if it's being used)
	if err := p.Lock(); err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	pushToken := fmt.Sprintf("%08X", rand.IntN(1<<(4*tokenLength)))
	var wg sync.WaitGroup
	wg.Add(1)
	process = &Process{
		wg:    &wg,
		token: pushToken,
	}
	go pushButton(process, isLong)

	// return token to get status of process
	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, pushToken)
	return

}

func PushStatus(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Push-Token")
	if len(token) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Return "Completed" if token isn't correct
	if process.token == token {
		io.WriteString(w, "Completed or token isn't correct")
		return
	}

	// Wait until the process is complete and return result
	process.wg.Wait()
	if process.err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, fmt.Sprintf("An error occured:%s", process.err))
	} else {
		io.WriteString(w, "Completed")
	}
}

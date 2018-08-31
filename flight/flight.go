//GO File
/*
 *Filename: flight/flight.go
 *
 *Author: kesheng, yuliang.cyl@alibaba-inc.com
 *Description: ---
 *Create: 2018-08-31 16:57:01
 *Last Modified: 2018-08-31 23:52:36
 */
package flight

import (
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Flight struct {
	end  chan int
	lock *sync.RWMutex

	isFighting bool

	ws      *websocket.Conn
	isLeave bool

	Name string
}

func NewFlight(ws *websocket.Conn) *Flight {

	f := &Flight{
		end:        make(chan int),
		lock:       &sync.RWMutex{},
		ws:         ws,
		isLeave:    false,
		isFighting: false,
	}

	buffer := make([]byte, 1024)
	n, err := ws.Read(buffer)
	if err != nil {
		return nil
	}

	f.Name = string(buffer[:n])

	log.Println(f.Name, "come!")

	return f
}

func (this *Flight) Notify(data []byte) error {
	sent := 0
	for sent < len(data) {
		n, err := this.ws.Write(data)
		sent += n
		if err != nil {
			log.Println(this.Name, "send data fail", err.Error())
			this.Leave()
			return err
		}
	}

	return nil
}

func (this *Flight) Read() ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := this.ws.Read(buffer)
	if err != nil {
		log.Println(this.Name, "read data fail", err.Error())
		this.Leave()
		return nil, err
	}

	return buffer[:n], err
}

func (this *Flight) Leave() {
	this.isLeave = true
	this.ws.Close()
	this.end <- 1
}

func (this *Flight) Wait() {
	<-this.end
}

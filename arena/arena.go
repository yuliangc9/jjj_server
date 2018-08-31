//GO File
/*
 *Filename: arena/arena.go
 *
 *Author: kesheng, yuliang.cyl@alibaba-inc.com
 *Description: ---
 *Create: 2018-08-31 16:37:21
 *Last Modified: 2018-08-31 23:12:05
 */
package arena

import (
	"log"
	"sync"

	"gitlab.alibaba-inc/jjj_server/battle"
	"gitlab.alibaba-inc/jjj_server/flight"
)

var Roma *Arena

type Arena struct {
	lock   *sync.Mutex
	waiter *flight.Flight
}

func init() {
	Roma = NewArena()
}

func NewArena() *Arena {
	return &Arena{
		lock:   &sync.Mutex{},
		waiter: nil,
	}
}

func (this *Arena) OnPlayer(f *flight.Flight) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.waiter == nil {
		log.Println(f.Name, "waiting...")
		this.waiter = f
		return
	}

	log.Println(f.Name, "match!", this.waiter.Name)
	b := battle.NewBattle(f, this.waiter)
	b.Fight()

	this.waiter = nil
}

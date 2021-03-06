//GO File
/*
 *Filename: battle/battle.go
 *
 *Author: kesheng, yuliang.cyl@alibaba-inc.com
 *Description: ---
 *Create: 2018-08-31 16:53:33
 *Last Modified: 2018-09-27 11:31:09
 */
package battle

import (
	"encoding/json"
	"log"

	"gitlab.alibaba-inc/jjj_server/flight"
)

var leaveInfo []byte
var beginSignal []byte

func init() {
	leaveInfo, _ = json.Marshal(map[string]bool{"leave": true})
}

type Battle struct {
	player1 *flight.Flight
	player2 *flight.Flight
}

func NewBattle(f1, f2 *flight.Flight) *Battle {
	return &Battle{
		player1: f1,
		player2: f2,
	}
}

func run(a, b *flight.Flight) {
	for {
		buffer, err := a.Read()
		if err != nil {
			log.Println(a.Name, "leave, notify", b.Name)
			b.Notify(leaveInfo)
			break
		}

		err = b.Notify(buffer)
		if err != nil {
			break
		}
	}
}

func (this *Battle) Fight() {
	go run(this.player1, this.player2)
	go run(this.player2, this.player1)
}

func (this *Battle) Check() error {
	//player1 is waiter
	player2Info, _ := json.Marshal(map[string]interface{}{
		"begin": true, "initHealth": this.player2.InitHealth, "name": this.player2.Name})
	err := this.player1.Notify(player2Info)
	if err != nil {
		return err
	}

	player1Info, _ := json.Marshal(map[string]interface{}{
		"begin": true, "initHealth": this.player1.InitHealth, "name": this.player1.Name})
	this.player2.Notify(player1Info)

	return nil
}

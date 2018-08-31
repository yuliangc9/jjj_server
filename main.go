package main

import (
	"net/http"

	"gitlab.alibaba-inc/jjj_server/arena"
	"gitlab.alibaba-inc/jjj_server/flight"
	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func FightServer(ws *websocket.Conn) {
	fly := flight.NewFlight(ws)
	if fly == nil {
		return
	}

	arena.Roma.OnPlayer(fly)

	fly.Wait()
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/fight", websocket.Handler(FightServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

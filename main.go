package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var base *websocket.Conn
var base1 *websocket.Conn

var leaveInfo []byte

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	nowTime := time.Now()
	fmt.Println("new client", nowTime.Unix())
	if base == nil {
		base = ws
	} else if base1 == nil {
		base1 = ws
	}

	buffer := make([]byte, 1024)
	for {
		n, err := ws.Read(buffer)
		if err != nil {
			if base == ws {
				base = nil
				if base1 != nil {
					base1.Write(leaveInfo)
				}
			} else if base1 == ws {
				base1 = nil
				if base != nil {
					base.Write(leaveInfo)
				}
			}

			fmt.Println("client leave", nowTime.Unix())
			break
		}

		if base != ws && base != nil {
			base.Write(buffer[:n])
		} else if base1 != ws && base1 != nil {
			base1.Write(buffer[:n])
		}
	}
}

func init() {
	leaveInfo, _ = json.Marshal(map[string]bool{"leave": true})
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/demo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

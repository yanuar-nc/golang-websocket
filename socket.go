package main

import (
	"time"

	"github.com/gorilla/websocket"
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(1)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	userCountTicker := time.NewTicker(filePeriod)

	defer func() {
		pingTicker.Stop()
		userCountTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-userCountTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readUserCount(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

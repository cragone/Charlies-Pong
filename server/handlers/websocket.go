package handlers

import (
	"game/tui"
	"game/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev only
}

func HandleSendGameStartUpWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer ws.Close()
	log.Println("client connected")

	msgs := []string{
		tui.WelcomingMessage,
		tui.GameCountDown3,
		tui.GameCountDown2,
		tui.GameCountDown1,
		tui.BeginingMessage,
	}

	const writeTimeout = 5 * time.Second

	for i, m := range msgs {
		select {
		case <-r.Context().Done():
			return
		default:
		}

		_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout))

		if err := ws.WriteJSON(utils.JSONResponse{"message": m}); err != nil {
			log.Printf("write json error: %v", err)
			return
		}
		if i < len(msgs)-1 {
			select {
			case <-time.After(1 * time.Second):
			case <-r.Context().Done():
				return
			}
		}
	}
}

package server

import (
	"fmt"
	"time"
)

type Message struct {
	room    *Room
	sender  string
	content string
	time    time.Time
	system  bool
}

func (s *Server) runBroadcast() {
	for msg := range s.messages {
		var formatted string
		if msg.system {
			formatted = fmt.Sprintf("[%s] 📢 %s %s\n",
				msg.time.Format("15:04:05"),
				msg.sender,
				msg.content,
			)
		} else {
			formatted = fmt.Sprintf("[%s][%s] %s\n",
				msg.time.Format("15:04:05"),
				msg.sender,
				msg.content,
			)
			fmt.Printf("[%s][%s][%s] %s\n",
				msg.time.Format("15:04:05"),
				msg.room.name,
				msg.sender,
				msg.content,
			)
		}
		msg.room.broadcast(formatted)
	}
}

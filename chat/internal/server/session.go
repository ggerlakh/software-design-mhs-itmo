package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	time.Sleep(50 * time.Millisecond)
	reader := bufio.NewReader(conn)

	for {
		room, nickname := s.showMainMenu(conn, reader)
		if room == nil {
			continue
		}

		client := &Client{conn: conn, nickname: nickname, room: room}
		count := room.addClient(client)

		fmt.Printf("[%s][%s] %s подключился (%d)\n",
			time.Now().Format("15:04:05"), room.name, nickname, count)

		send(conn, PrefixChat+"\n")
		send(conn, fmt.Sprintf("\n✅ Вы в комнате '%s' как %s\n", room.name, nickname))
		send(conn, "Команды: /quit, /users, /menu\n")
		send(conn, "─────────────────────────────────────\n")

		s.messages <- Message{
			room:    room,
			sender:  nickname,
			content: "присоединился к комнате",
			time:    time.Now(),
			system:  true,
		}

		backToMenu := s.handleChatLoop(client, reader)

		s.removeClientFromRoom(client)

		if !backToMenu {
			return
		}
	}
}

func (s *Server) handleChatLoop(c *Client, reader *bufio.Reader) bool {
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		if strings.HasPrefix(message, "/") {
			action := s.handleCommand(c, message)
			if action == ActionQuit {
				return false
			}
			if action == ActionMenu {
				return true
			}
			continue
		}

		s.messages <- Message{
			room:    c.room,
			sender:  c.nickname,
			content: message,
			time:    time.Now(),
			system:  false,
		}
	}
}

func (s *Server) removeClientFromRoom(c *Client) {
	count := c.room.removeClient(c)

	fmt.Printf("[%s][%s] %s отключился (%d)\n",
		time.Now().Format("15:04:05"), c.room.name, c.nickname, count)

	s.messages <- Message{
		room:    c.room,
		sender:  c.nickname,
		content: "покинул комнату",
		time:    time.Now(),
		system:  true,
	}

	if count == 0 {
		s.deleteRoom(c.room.name)
		fmt.Printf("[%s] Комната '%s' удалена\n",
			time.Now().Format("15:04:05"), c.room.name)
	}
}

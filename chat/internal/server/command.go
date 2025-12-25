package server

import (
	"fmt"
	"strings"
)

type CommandAction string

const (
	ActionNone CommandAction = ""
	ActionQuit CommandAction = "quit"
	ActionMenu CommandAction = "menu"
)

func (s *Server) handleCommand(c *Client, cmd string) CommandAction {
	switch cmd {
	case "/quit":
		send(c.conn, "До свидания!\n")
		return ActionQuit

	case "/menu":
		send(c.conn, "🔙 Возврат в главное меню...\n")
		return ActionMenu

	case "/users":
		users := c.room.listUsers()
		send(c.conn, fmt.Sprintf("👥 В комнате (%d): %s\n",
			len(users), strings.Join(users, ", ")))

	default:
		send(c.conn, "❓ Команды: /quit, /users, /menu\n")
	}
	return ActionNone
}

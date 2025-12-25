package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func (s *Server) showMainMenu(conn net.Conn, reader *bufio.Reader) (*Room, string) {
	send(conn, "\n╔════════════════════════════════════╗\n")
	send(conn, "║         🏠 ГЛАВНОЕ МЕНЮ            ║\n")
	send(conn, "╠════════════════════════════════════╣\n")
	send(conn, "║  1. Создать комнату                ║\n")
	send(conn, "║  2. Присоединиться к комнате       ║\n")
	send(conn, "╚════════════════════════════════════╝\n")

	for {
		choice, err := askInput(conn, reader, "Введите 1 или 2")
		if err != nil {
			return nil, ""
		}

		switch choice {
		case "1":
			room, nick := s.handleCreateRoom(conn, reader)
			if room == nil {
				return s.showMainMenu(conn, reader)
			}
			return room, nick
		case "2":
			room, nick := s.handleJoinRoom(conn, reader)
			if room == nil {
				return s.showMainMenu(conn, reader)
			}
			return room, nick
		default:
			send(conn, "❌ Неверный выбор.\n")
		}
	}
}

func (s *Server) handleCreateRoom(conn net.Conn, reader *bufio.Reader) (*Room, string) {
	send(conn, "\n📝 СОЗДАНИЕ КОМНАТЫ\n")
	send(conn, "(введите /menu для возврата)\n\n")

	for {
		roomName, err := askInput(conn, reader, "Введите название комнаты")
		if err != nil {
			return nil, ""
		}
		if roomName == "/menu" {
			return nil, ""
		}
		if roomName == "" {
			send(conn, "❌ Название не может быть пустым\n")
			continue
		}

		if s.roomExists(roomName) {
			send(conn, "❌ Комната уже существует\n")
			continue
		}

		password, err := askInput(conn, reader, "Введите пароль для комнаты")
		if err != nil {
			return nil, ""
		}
		if password == "/menu" {
			return nil, ""
		}

		nickname, err := askInput(conn, reader, "Введите ваш никнейм")
		if err != nil {
			return nil, ""
		}
		if nickname == "/menu" {
			return nil, ""
		}
		if nickname == "" {
			nickname = fmt.Sprintf("User_%d", time.Now().UnixNano()%10000)
		}

		room := s.createRoom(roomName, password)
		fmt.Printf("[%s] Комната '%s' создана\n", time.Now().Format("15:04:05"), roomName)
		return room, nickname
	}
}

func (s *Server) handleJoinRoom(conn net.Conn, reader *bufio.Reader) (*Room, string) {
	send(conn, "\n🚪 ВХОД В КОМНАТУ\n")
	send(conn, "(введите /menu для возврата)\n\n")

	var room *Room
	for {
		roomName, err := askInput(conn, reader, "Введите название комнаты")
		if err != nil {
			return nil, ""
		}
		if roomName == "/menu" {
			return nil, ""
		}
		if roomName == "" {
			send(conn, "❌ Название не может быть пустым\n")
			continue
		}

		room = s.getRoom(roomName)
		if room == nil {
			send(conn, "❌ Комната не найдена\n")
			continue
		}
		break
	}

	for attempts := 0; attempts < 3; attempts++ {
		remaining := 3 - attempts
		password, err := askInput(conn, reader, fmt.Sprintf("Введите пароль (попыток: %d)", remaining))
		if err != nil {
			return nil, ""
		}
		if password == "/menu" {
			return nil, ""
		}

		if room.checkPassword(password) {
			nickname, err := askInput(conn, reader, "Введите ваш никнейм")
			if err != nil {
				return nil, ""
			}
			if nickname == "/menu" {
				return nil, ""
			}
			if nickname == "" {
				nickname = fmt.Sprintf("User_%d", time.Now().UnixNano()%10000)
			}
			return room, nickname
		}

		if attempts < 2 {
			send(conn, "❌ Неверный пароль\n")
		}
	}

	send(conn, "❌ Превышено количество попыток\n")
	return nil, ""
}

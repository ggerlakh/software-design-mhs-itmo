package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	port     string
	rooms    map[string]*Room
	roomsMu  sync.RWMutex
	messages chan Message
}

func New(port string) *Server {
	return &Server{
		port:     port,
		rooms:    make(map[string]*Room),
		messages: make(chan Message, 100),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp4", ":"+s.port)
	if err != nil {
		return fmt.Errorf("не удалось запустить сервер: %w", err)
	}
	defer listener.Close()

	fmt.Printf("🚀 Сервер запущен на порту %s\n", s.port)
	fmt.Println("Ожидание подключений...")

	go s.runBroadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) roomExists(name string) bool {
	s.roomsMu.RLock()
	defer s.roomsMu.RUnlock()
	_, exists := s.rooms[name]
	return exists
}

func (s *Server) getRoom(name string) *Room {
	s.roomsMu.RLock()
	defer s.roomsMu.RUnlock()
	return s.rooms[name]
}

func (s *Server) createRoom(name, password string) *Room {
	room := newRoom(name, password)
	s.roomsMu.Lock()
	s.rooms[name] = room
	s.roomsMu.Unlock()
	return room
}

func (s *Server) deleteRoom(name string) {
	s.roomsMu.Lock()
	delete(s.rooms, name)
	s.roomsMu.Unlock()
}

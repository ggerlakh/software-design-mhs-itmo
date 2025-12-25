package server

import (
	"net"
	"sync"
)

type Room struct {
	name         string
	passwordHash string
	clients      map[net.Conn]*Client
	mu           sync.RWMutex
}

type Client struct {
	conn     net.Conn
	nickname string
	room     *Room
}

func newRoom(name, password string) *Room {
	return &Room{
		name:         name,
		passwordHash: hashPassword(password),
		clients:      make(map[net.Conn]*Client),
	}
}

func (r *Room) addClient(c *Client) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c.conn] = c
	return len(r.clients)
}

func (r *Room) removeClient(c *Client) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c.conn)
	return len(r.clients)
}

func (r *Room) listUsers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := make([]string, 0, len(r.clients))
	for _, c := range r.clients {
		users = append(users, c.nickname)
	}
	return users
}

func (r *Room) broadcast(msg string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.clients {
		send(c.conn, msg)
	}
}

func (r *Room) checkPassword(password string) bool {
	return hashPassword(password) == r.passwordHash
}

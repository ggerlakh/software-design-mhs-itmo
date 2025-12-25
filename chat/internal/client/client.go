package client

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	addr     string
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	chatMode bool
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) connect() error {
	conn, err := net.Dial("tcp4", c.addr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться: %w", err)
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)

	return nil
}

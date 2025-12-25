package server

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"
)

const (
	PrefixPrompt = "PROMPT:"
	PrefixChat   = "CHAT:"
)

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func send(conn net.Conn, msg string) {
	conn.Write([]byte(msg))
}

func askInput(conn net.Conn, reader *bufio.Reader, prompt string) (string, error) {
	conn.Write([]byte(PrefixPrompt + prompt + "\n"))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

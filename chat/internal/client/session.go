package client

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func (c *Client) Start() error {
	if err := c.connect(); err != nil {
		return err
	}
	defer c.conn.Close()

	fmt.Printf("🔗 Подключено к %s\n", c.addr)

	done := make(chan struct{})
	stdinReader := bufio.NewReader(os.Stdin)

	go c.handleSignals(done)

	for {
		select {
		case <-done:
			return nil
		default:
		}

		message, err := c.reader.ReadString('\n')
		if err != nil {
			select {
			case <-done:
				return nil
			default:
				fmt.Println("\n❌ Соединение потеряно")
				return nil
			}
		}

		message = strings.TrimSuffix(message, "\n")
		c.processMessage(message, done, stdinReader)
	}
}

func (c *Client) handleSignals(done chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\n👋 Выход...")
	close(done)
	c.conn.Close()
}

func (c *Client) processMessage(message string, done chan struct{}, stdinReader *bufio.Reader) {
	if strings.HasPrefix(message, PrefixPrompt) {
		promptText := strings.TrimPrefix(message, PrefixPrompt)
		fmt.Print(promptText + ": ")

		input, err := stdinReader.ReadString('\n')
		if err != nil {
			return
		}

		c.writer.WriteString(input)
		c.writer.Flush()

	} else if strings.HasPrefix(message, PrefixChat) {
		c.chatMode = true
		go c.runChatMode(done, stdinReader)

	} else {
		fmt.Println(message)
	}
}

func (c *Client) runChatMode(done chan struct{}, stdinReader *bufio.Reader) {
	for {
		select {
		case <-done:
			return
		default:
		}

		input, err := stdinReader.ReadString('\n')
		if err != nil {
			return
		}

		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			continue
		}

		if !strings.HasPrefix(trimmed, "/") {
			fmt.Print(CursorUp + ClearLine + CursorStart)
		}

		c.writer.WriteString(input)
		c.writer.Flush()

		if trimmed == "/quit" {
			close(done)
			return
		}
		if trimmed == "/menu" {
			c.chatMode = false
			return
		}
	}
}

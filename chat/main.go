package main

import (
	"flag"
	"fmt"
	"os"

	"chat/internal/client"
	"chat/internal/server"
)

func main() {
	mode := flag.String("mode", "", "Режим: server или client")
	port := flag.String("port", "8080", "Порт сервера")
	host := flag.String("host", "localhost", "Адрес сервера (для клиента)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
💬 CHAT — чат с комнатами

Использование:
  chat -mode=server [-port=8080]
  chat -mode=client -host=IP [-port=8080]

Примеры:
  chat -mode=server -port=9000
  chat -mode=client -host=192.168.1.100 -port=9000

`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch *mode {
	case "server":
		runServer(*port)
	case "client":
		runClient(*host, *port)
	default:
		fmt.Fprintf(os.Stderr, "❌ Неизвестный режим: %s\n", *mode)
		os.Exit(1)
	}
}

func runServer(port string) {
	fmt.Println(`
╔══════════════════════════════════════╗
║          🖥️  РЕЖИМ СЕРВЕРА           ║
╚══════════════════════════════════════╝`)

	srv := server.New(port)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func runClient(host, port string) {
	fmt.Println(`
╔══════════════════════════════════════╗
║          💻 РЕЖИМ КЛИЕНТА            ║
╚══════════════════════════════════════╝`)

	addr := fmt.Sprintf("%s:%s", host, port)
	cli := client.New(addr)

	if err := cli.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
		os.Exit(1)
	}
}

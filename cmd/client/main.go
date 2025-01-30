package main

import (
	"flag"
	"fmt"
	"net"
	"public-domain/internal/client"
)

func main() {
	serverAddr := flag.String("server-addr", "127.0.0.1", "Адрес сервера")
	remotePort := flag.String("remote-port", "8080", "Порт сервера для входящих соединений")
	localPort := flag.String("local-port", "3000", "Локальный порт, который нужно расшарить")
	flag.Parse()

	conn, err := net.Dial("tcp", *serverAddr+":"+*remotePort)
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу", *serverAddr, "на порту", *remotePort)

	for {
		client.HandleTunnel(conn, *localPort)
	}
}

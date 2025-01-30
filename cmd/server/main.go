package main

import (
	"flag"
	"fmt"
	"net"
	"public-domain/internal/server"
)

func main() {
	port := flag.String("port", "8080", "Порт, на котором сервер будет принимать подключения")
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту", *port)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		fmt.Println("Подключился клиент", clientConn.RemoteAddr())

		// Ждём соединение от клиента для проброса
		targetConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при ожидании соединения от клиента:", err)
			clientConn.Close()
			continue
		}

		go server.HandleClientConnection(clientConn, targetConn)
	}
}

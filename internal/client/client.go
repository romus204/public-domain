package client

import (
	"fmt"
	"io"
	"net"
	"sync"
)

func HandleTunnel(serverConn net.Conn, localPort string) {
	defer serverConn.Close()

	localConn, err := net.Dial("tcp", "127.0.0.1:"+localPort)
	if err != nil {
		fmt.Println("Ошибка соединения с локальным сервером:", err)
		return
	}
	defer localConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Перенаправление данных в обе стороны
	go func() {
		defer wg.Done()
		io.Copy(localConn, serverConn)
	}()
	go func() {
		defer wg.Done()
		io.Copy(serverConn, localConn)
	}()

	wg.Wait()
}

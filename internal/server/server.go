package server

import (
	"io"
	"net"
	"sync"
)

func HandleClientConnection(clientConn net.Conn, targetConn net.Conn) {
	defer clientConn.Close()
	defer targetConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// Перенаправление данных в обе стороны
	go func() {
		defer wg.Done()
		io.Copy(targetConn, clientConn)
	}()
	go func() {
		defer wg.Done()
		io.Copy(clientConn, targetConn)
	}()

	wg.Wait()
}

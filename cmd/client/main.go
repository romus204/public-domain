package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func main() {
	serverAddr := flag.String("server-addr", "127.0.0.1", "Адрес сервера")
	serverPort := flag.String("server-port", "8383", "Порт сервера")
	localPort := flag.String("local-port", "3002", "Локальный порт, который нужно расшарить")
	caFile := flag.String("ca", "", "Файл корневого сертификата (опционально)")
	flag.Parse()

	// Настройка конфигурации TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Проверять сертификат сервера
	}

	if *caFile != "" {
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatalf("Ошибка загрузки корневого сертификата: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	var prevServerConn net.Conn

	first := true

	for {
		// Подключение к серверу
		serverConn, err := tls.Dial("tcp", *serverAddr+":"+*serverPort, tlsConfig)
		if !first {
			prevServerConn.Close()
		}
		prevServerConn = serverConn
		if err != nil {
			log.Println("Ошибка подключения к серверу:", err)
			time.Sleep(1 * time.Second) // Пауза перед повторной попыткой
			continue
		}

		fmt.Println("Подключено к серверу", *serverAddr, "на порту", *serverPort)

		// Подключение к локальному сервису
		localConn, err := net.Dial("tcp", ":"+*localPort)
		if err != nil {
			log.Println("Ошибка подключения к локальному сервису:", err)
			serverConn.Close()
			time.Sleep(1 * time.Second) // Пауза перед повторной попыткой
			continue
		}

		fmt.Println("Подключено к локальному сервису на порту", *localPort)

		// Пересылка данных между сервером и локальным сервисом
		go func() {
			_, err := io.Copy(localConn, serverConn)
			if err != nil {
				log.Println("Ошибка при пересылке данных от сервера к локальному сервису:", err)
			}
		}()

		go func() {
			_, err := io.Copy(serverConn, localConn)
			if err != nil {
				log.Println("Ошибка при пересылке данных от локального сервиса к серверу:", err)
			}
		}()

		first = false
	}
}

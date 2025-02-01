package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	port := flag.String("port", "8383", "Порт, на котором сервер будет принимать подключения")
	certFile := flag.String("cert", "server.crt", "Файл сертификата")
	keyFile := flag.String("key", "server.key", "Файл приватного ключа")
	caFile := flag.String("ca", "", "Файл корневого сертификата (опционально)")
	flag.Parse()

	// Загрузка сертификата и приватного ключа
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки сертификата и ключа: %v", err)
	}

	// Настройка конфигурации TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert, // или tls.RequireAndVerifyClientCert для проверки клиентских сертификатов
	}

	if *caFile != "" {
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			log.Fatalf("Ошибка загрузки корневого сертификата: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	// Создание TLS-листенера
	listener, err := tls.Listen("tcp", ":"+*port, tlsConfig)
	if err != nil {
		log.Fatalf("Ошибка запуска TLS-сервера: %v", err)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту", *port)

	for {
		// Ожидание подключения клиента (локального сервиса)
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка при принятии соединения от клиента:", err)
			continue
		}

		fmt.Println("Клиент подключен:", clientConn.RemoteAddr())

		// Ожидание подключения внешнего пользователя
		userConn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка при принятии соединения от пользователя:", err)
			clientConn.Close()
			continue
		}

		fmt.Println("Пользователь подключен:", userConn.RemoteAddr())

		// Пересылка данных между клиентом и пользователем
		go func() {
			_, err := io.Copy(clientConn, userConn)
			if err != nil {
				log.Println("Ошибка при пересылке данных от пользователя к клиенту:", err)
			}
			clientConn.Close()
			userConn.Close()
		}()

		go func() {
			_, err := io.Copy(userConn, clientConn)
			if err != nil {
				log.Println("Ошибка при пересылке данных от клиента к пользователю:", err)
			}
			clientConn.Close()
			userConn.Close()
		}()
	}
}

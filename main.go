package main

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel

	addr := "0.0.0.0:8080"

	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Не удалось открыть соединение: %v", err)
	}
	defer l.Close()

	logger.Infof("Запущен сервер %s...", addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorf("Ошибка соединения: %v", err)
		} else {
			go handleRecvPkg(conn, logger)
		}
	}
}

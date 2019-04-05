package main

import (
	"net"

	"github.com/nburunova/go-egts/egts"
	"github.com/sirupsen/logrus"
)

func handleRecvPkg(conn net.Conn, logger *logrus.Logger) {
	buf := make([]byte, 1024)

	logger.Warnf("Установлено соединение с %s", conn.RemoteAddr())

	for {

		_, err := conn.Read(buf)
		if err != nil {
			logger.Debug("Не смогли принять пакет")
		}

		logger.Debugf("Принят пакет: %X\v", buf)
		//printDecodePackage(buf)

		resp, err := egts.ParsePacket(buf)

		conn.Write(resp)

	}
}

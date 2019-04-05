package main

import (
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/nburunova/go-egts/egts"
	"github.com/sirupsen/logrus"
)

const timerDuration = 5 * time.Second

func handleRecvPkg(conn net.Conn, logger *logrus.Logger) {
	var (
		srResultCodePkg   []byte
		serviceType       uint8
		srResponsesRecord egts.RecordDataSet
	)
	buf := make([]byte, 1024)

	logger.Warnf("Установлено соединение с %s", conn.RemoteAddr())

	for {
	Received:
		serviceType = 0
		srResponsesRecord = nil
		srResultCodePkg = nil

		pkgLen, err := conn.Read(buf)

		connTimer := time.NewTimer(timerDuration)
		switch err {
		case nil:
			connTimer.Reset(timerDuration)
			logger.Debugf("Принят пакет: %X\v", buf[:pkgLen])
			break
		case io.EOF:
			<-connTimer.C
			conn.Close()
			logger.Warnf("Соединение %s закрыто по таймауту", conn.RemoteAddr())
			return
		default:
			logger.Errorf("Ошибка при получении:", err)
			conn.Close()
			return
		}

		logger.Debugf("Принят пакет: %X\v", buf)
		//printDecodePackage(buf)

		pkg := egts.EgtsPackage{}
		resultCode, err := pkg.Decode(buf[:pkgLen])
		if err != nil {
			logger.Errorf("Не удалось расшифровать пакет: %v", err)

			resp, err := pkg.CreatePtResponse(resultCode, serviceType, nil)
			if err != nil {
				logger.Errorf("Ошибка сборки ответа EGTS_PT_RESPONSE с ошибкой: %v", err)
				goto Received
			}
			conn.Write(resp)

			//printDecodePackage("Отправлен пакет EGTS_PT_RESPONSE", resp)
			goto Received
		}

		switch pkg.PacketType {
		case egts.EgtsPtAppdataPkgType:
			logger.Info("Тип пакета EGTS_PT_APPDATA")

			for _, rec := range *pkg.ServicesFrameData.(*egts.ServiceDataSet) {
				exportPacket := egts.EgtsParsePacket{
					PacketID: uint32(pkg.PacketIdentifier),
				}
				packetIdBytes := make([]byte, 4)

				srResponsesRecord = append(srResponsesRecord, egts.RecordData{
					SubrecordType:   egts.EgtsSrRecordResponsePkgType,
					SubrecordLength: 3,
					SubrecordData: &egts.EgtsSrResponse{
						ConfirmedRecordNumber: rec.RecordNumber,
						RecordStatus:          egts.EgtsPcOk,
					},
				})
				serviceType = rec.SourceServiceType
				logger.Info("Тип сервиса ", serviceType)

				exportPacket.Client = rec.ObjectIdentifier

				for _, subRec := range rec.RecordDataSet {
					switch subRecData := subRec.SubrecordData.(type) {
					case *egts.EgtsSrTermIdentity:
						logger.Debugf("Разбор подзаписи EGTS_SR_TERM_IDENTITY")
						if srResultCodePkg, err = pkg.CreateSrResultCode(egts.EgtsPcOk); err != nil {
							logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
						}
					case *egts.EgtsSrAuthInfo:
						logger.Debugf("Разбор подзаписи EGTS_SR_AUTH_INFO")
						if srResultCodePkg, err = pkg.CreateSrResultCode(egts.EgtsPcOk); err != nil {
							logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
						}
					case *egts.EgtsSrResponse:
						logger.Debugf("Разбор подзаписи EGTS_SR_RESPONSE")
						goto Received
					case *egts.EgtsSrPosData:
						logger.Debugf("Разбор подзаписи EGTS_SR_POS_DATA")
						exportPacket.NavigationTime = subRecData.NavigationTime
						exportPacket.Latitude = subRecData.Latitude
						exportPacket.Longitude = subRecData.Longitude
						exportPacket.Speed = subRecData.Speed
						exportPacket.Course = subRecData.Direction
					case *egts.EgtsSrExtPosData:
						logger.Debugf("Разбор подзаписи EGTS_SR_EXT_POS_DATA")
						exportPacket.Nsat = subRecData.Satellites
						exportPacket.Pdop = subRecData.PositionDilutionOfPrecision

					case *egts.EgtsSrAdSensorsData:
						logger.Debugf("Разбор подзаписи EGTS_SR_AD_SENSORS_DATA")

						exportPacket.AnSensors = make(map[uint8]uint32)
						exportPacket.AnSensors[1] = subRecData.AnalogSensor1
						exportPacket.AnSensors[2] = subRecData.AnalogSensor2
						exportPacket.AnSensors[3] = subRecData.AnalogSensor3
						exportPacket.AnSensors[4] = subRecData.AnalogSensor4
						exportPacket.AnSensors[5] = subRecData.AnalogSensor5
						exportPacket.AnSensors[6] = subRecData.AnalogSensor6
						exportPacket.AnSensors[7] = subRecData.AnalogSensor7
						exportPacket.AnSensors[8] = subRecData.AnalogSensor8
					case *egts.EgtsSrAbsCntrData:
						logger.Debugf("Разбор подзаписи EGTS_SR_ABS_CNTR_DATA")

						switch subRecData.CounterNumber {
						case 110:
							// Три младших байта номера передаваемой записи (идет вместе с каждой POS_DATA).
							binary.BigEndian.PutUint32(packetIdBytes, subRecData.CounterValue)
							exportPacket.PacketID = subRecData.CounterValue
						case 111:
							// один старший байт номера передаваемой записи (идет вместе с каждой POS_DATA).
							tmpBuf := make([]byte, 4)
							binary.BigEndian.PutUint32(tmpBuf, subRecData.CounterValue)

							if len(packetIdBytes) == 4 {
								packetIdBytes[3] = tmpBuf[3]
							} else {
								packetIdBytes = tmpBuf
							}

							exportPacket.PacketID = binary.LittleEndian.Uint32(packetIdBytes)
						}
					case *egts.EgtsSrLiquidLevelSensor:
						logger.Debugf("Разбор подзаписи EGTS_SR_LIQUID_LEVEL_SENSOR")
						sensorData := egts.LiquidSensor{
							SensorNumber: subRecData.LiquidLevelSensorNumber,
							ErrorFlag:    subRecData.LiquidLevelSensorErrorFlag,
						}

						switch subRecData.LiquidLevelSensorValueUnit {
						case "00", "01":
							sensorData.ValueMm = subRecData.LiquidLevelSensorData
						case "10":
							sensorData.ValueL = subRecData.LiquidLevelSensorData * 10
						}

						exportPacket.LiquidSensors = append(exportPacket.LiquidSensors, sensorData)
					}
				}
			}

			resp, err := pkg.CreatePtResponse(resultCode, serviceType, srResponsesRecord)
			if err != nil {
				logger.Errorf("Ошибка сборки ответа: %v", err)
				goto Received
			}
			conn.Write(resp)

			logger.Debugf("Отправлен пакет EGTS_PT_RESPONSE: %X", resp)
			//logger.Debug(printDecodePackage(resp))

			if len(srResultCodePkg) > 0 {
				conn.Write(srResultCodePkg)
				logger.Debugf("Отправлен пакет EGTS_SR_RESULT_CODE: %X", resp)
				//logger.Debug(printDecodePackage(srResultCodePkg))
			}
		case egts.EgtsPtResponsePkgType:
			logger.Printf("Тип пакета EGTS_PT_RESPONSE")
		}

	}
}

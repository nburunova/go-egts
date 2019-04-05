package egts

import (
	"encoding/binary"

	"github.com/sirupsen/logrus"
)

func ParsePacket(buf []byte, logger *logrus.Logger) ([]byte, error) {
	var (
		srResultCodePkg   []byte
		serviceType       uint8
		srResponsesRecord RecordDataSet
	)

	pkgLen := len(buf)
	pkg := EgtsPackage{}
	resultCode, err := pkg.Decode(buf[:pkgLen])

	switch pkg.PacketType {
	case egtsPtAppdata:
		logger.Info("Тип пакета EGTS_PT_APPDATA")

		for _, rec := range *pkg.ServicesFrameData.(*ServiceDataSet) {
			exportPacket := EgtsParsePacket{
				PacketID: uint32(pkg.PacketIdentifier),
			}
			packetIdBytes := make([]byte, 4)

			srResponsesRecord = append(srResponsesRecord, RecordData{
				SubrecordType:   egtsSrRecordResponse,
				SubrecordLength: 3,
				SubrecordData: &EgtsSrResponse{
					ConfirmedRecordNumber: rec.RecordNumber,
					RecordStatus:          egtsPcOk,
				},
			})
			serviceType = rec.SourceServiceType
			logger.Info("Тип сервиса ", serviceType)

			exportPacket.Client = rec.ObjectIdentifier

			for _, subRec := range rec.RecordDataSet {
				switch subRecData := subRec.SubrecordData.(type) {
				case *EgtsSrTermIdentity:
					logger.Debugf("Разбор подзаписи EGTS_SR_TERM_IDENTITY")
					if srResultCodePkg, err = pkg.CreateSrResultCode(egtsPcOk); err != nil {
						logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
					}
				case *EgtsSrAuthInfo:
					logger.Debugf("Разбор подзаписи EGTS_SR_AUTH_INFO")
					if srResultCodePkg, err = pkg.CreateSrResultCode(egtsPcOk); err != nil {
						logger.Errorf("Ошибка сборки EGTS_SR_RESULT_CODE: %v", err)
					}
				case *EgtsSrPosData:
					logger.Debugf("Разбор подзаписи EGTS_SR_POS_DATA")

					exportPacket.NavigationTime = subRecData.NavigationTime
					exportPacket.Latitude = subRecData.Latitude
					exportPacket.Longitude = subRecData.Longitude
					exportPacket.Speed = subRecData.Speed
					exportPacket.Course = subRecData.Direction
				case *EgtsSrExtPosData:
					logger.Debugf("Разбор подзаписи EGTS_SR_EXT_POS_DATA")
					exportPacket.Nsat = subRecData.Satellites
					exportPacket.Pdop = subRecData.PositionDilutionOfPrecision

				case *EgtsSrAdSensorsData:
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
				case *EgtsSrAbsCntrData:
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
				case *EgtsSrLiquidLevelSensor:
					logger.Debugf("Разбор подзаписи EGTS_SR_LIQUID_LEVEL_SENSOR")
					sensorData := LiquidSensor{
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

		logger.Debugf("Отправлен пакет EGTS_PT_RESPONSE: %X", resp)
		//logger.Debug(printDecodePackage(resp))

		if len(srResultCodePkg) > 0 {
			logger.Debugf("Отправлен пакет EGTS_SR_RESULT_CODE: %X", resp)
			//logger.Debug(printDecodePackage(srResultCodePkg))
		}
		return resp, err
		// case egtsPtResponse:
		// 	logger.Printf("Тип пакета EGTS_PT_RESPONSE")
	}
	return make([]byte, 0), err
}

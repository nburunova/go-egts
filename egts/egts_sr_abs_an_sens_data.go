package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type EgtsSrAbsAnSensData struct {
	AnalogSensorNumber uint8  `json:"ASN"`
	AnalogSensorValue  uint32 `json:"ASV"`
}

func (e *EgtsSrAbsAnSensData) Decode(content []byte) error {
	var (
		err error
	)
	buf := bytes.NewReader(content)

	if e.AnalogSensorNumber, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить номер аналогового входа: %v", err)
	}

	tmpBuf := make([]byte, 3)
	if _, err = buf.Read(tmpBuf); err != nil {
		return fmt.Errorf("Не удалось получить значение показаний аналогового входа: %v", err)
	}

	ansensVal := append(tmpBuf, 0x00)
	e.AnalogSensorValue = binary.LittleEndian.Uint32(ansensVal)

	return err
}

func (e *EgtsSrAbsAnSensData) Encode() ([]byte, error) {
	var (
		err    error
		result []byte
	)
	buf := new(bytes.Buffer)

	if err = buf.WriteByte(e.AnalogSensorNumber); err != nil {
		return result, fmt.Errorf("Не удалось записать номер аналогового входа: %v", err)
	}

	counterVal := make([]byte, 4)
	binary.LittleEndian.PutUint32(counterVal, e.AnalogSensorValue)
	if _, err = buf.Write(counterVal[:3]); err != nil {
		return result, fmt.Errorf("Не удалось записать значение показаний аналогового входа: %v", err)
	}

	result = buf.Bytes()
	return result, err
}

func (e *EgtsSrAbsAnSensData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

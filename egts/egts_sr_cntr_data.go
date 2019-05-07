package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type EgtsSrCntrsData struct {
	CounterFieldExists1 string `json:"CFE1"`
	CounterFieldExists2 string `json:"CFE2"`
	CounterFieldExists3 string `json:"CFE3"`
	CounterFieldExists4 string `json:"CFE4"`
	CounterFieldExists5 string `json:"CFE5"`
	CounterFieldExists6 string `json:"CFE6"`
	CounterFieldExists7 string `json:"CFE7"`
	CounterFieldExists8 string `json:"CFE8"`
	Counter1            uint32 `json:"ANS1"`
	Counter2            uint32 `json:"ANS2"`
	Counter3            uint32 `json:"ANS3"`
	Counter4            uint32 `json:"ANS4"`
	Counter5            uint32 `json:"ANS5"`
	Counter6            uint32 `json:"ANS6"`
	Counter7            uint32 `json:"ANS7"`
	Counter8            uint32 `json:"ANS8"`
}

func (e *EgtsSrCntrsData) Decode(content []byte) error {
	var (
		err        error
		flags      byte
		counterVal []byte
	)
	buf := bytes.NewReader(content)

	//байт флагов
	if flags, err = buf.ReadByte(); err != nil {
		return fmt.Errorf("Не удалось получить байт цифровых выходов ad_sesor_data: %v", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)

	e.CounterFieldExists1 = flagBits[:1]
	e.CounterFieldExists2 = flagBits[1:2]
	e.CounterFieldExists3 = flagBits[2:3]
	e.CounterFieldExists4 = flagBits[3:4]
	e.CounterFieldExists5 = flagBits[4:5]
	e.CounterFieldExists6 = flagBits[5:6]
	e.CounterFieldExists7 = flagBits[6:7]
	e.CounterFieldExists8 = flagBits[7:]

	tmpBuf := make([]byte, 3)
	if e.CounterFieldExists1 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS1: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter1 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists2 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS2: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter2 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists3 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS3: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter3 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists4 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS4: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter4 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists5 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS5: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter5 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists6 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS6: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter6 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists7 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS7: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter7 = binary.LittleEndian.Uint32(counterVal)
	}

	if e.CounterFieldExists8 == "1" {
		if _, err = buf.Read(tmpBuf); err != nil {
			return fmt.Errorf("Не удалось получить показания ANS8: %v", err)
		}
		counterVal = append(tmpBuf, 0x00)
		e.Counter8 = binary.LittleEndian.Uint32(counterVal)
	}
	return err
}

func (e *EgtsSrCntrsData) Encode() ([]byte, error) {
	var (
		err    error
		flags  uint64
		result []byte
	)

	buf := new(bytes.Buffer)

	flagsBits := e.CounterFieldExists1 +
		e.CounterFieldExists2 +
		e.CounterFieldExists3 +
		e.CounterFieldExists4 +
		e.CounterFieldExists5 +
		e.CounterFieldExists6 +
		e.CounterFieldExists7 +
		e.CounterFieldExists8

	if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
		return result, fmt.Errorf("Не удалось сгенерировать байт байт цифровых выходов ad_sesor_data: %v", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("Не удалось записать байт флагов ext_pos_data: %v", err)
	}

	cntrVal := make([]byte, 4)
	if e.CounterFieldExists1 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter1)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS1: %v", err)
		}
	}

	if e.CounterFieldExists2 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter2)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS2: %v", err)
		}
	}

	if e.CounterFieldExists3 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter3)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS3: %v", err)
		}
	}

	if e.CounterFieldExists4 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter4)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS4: %v", err)
		}
	}

	if e.CounterFieldExists5 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter5)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS5: %v", err)
		}
	}

	if e.CounterFieldExists6 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter6)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS6: %v", err)
		}
	}

	if e.CounterFieldExists7 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter7)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS7: %v", err)
		}
	}

	if e.CounterFieldExists8 == "1" {
		binary.LittleEndian.PutUint32(cntrVal, e.Counter8)
		if _, err = buf.Write(cntrVal[:3]); err != nil {
			return result, fmt.Errorf("Не удалось запистаь показания ANS8: %v", err)
		}
	}

	result = buf.Bytes()
	return result, err
}

func (e *EgtsSrCntrsData) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

package egts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//EgtsSrDispatcherIdentity --
type EgtsSrDispatcherIdentity struct {
	DispatcherType uint8  //DT
	DispatcherID   uint32 //DID
	Description    string // DSCR
}

func (e *EgtsSrDispatcherIdentity) Decode(content []byte) error {
	i := 0
	e.DispatcherType = uint8(content[0])
	i++
	e.DispatcherID = binary.LittleEndian.Uint32(content[i : i+4])
	i += 4
	e.Description = string(content[i:len(content)])
	return nil
}

func (e *EgtsSrDispatcherIdentity) Encode() ([]byte, error) {
	var (
		result []byte
		err    error
	)
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, e.DispatcherType); err != nil {
		return result, fmt.Errorf("Не удалось записать тип диспетчера")
	}

	if err = binary.Write(buf, binary.LittleEndian, e.DispatcherID); err != nil {
		return result, fmt.Errorf("Не удалось записать ID диспетчера")
	}

	if err = binary.Write(buf, binary.LittleEndian, e.Description); err != nil {
		return result, fmt.Errorf("Не удалось записать описание диспетчера")
	}
	result = buf.Bytes()
	return result, err
}

func (e *EgtsSrDispatcherIdentity) Length() uint16 {
	var result uint16

	if recBytes, err := e.Encode(); err != nil {
		result = uint16(0)
	} else {
		result = uint16(len(recBytes))
	}

	return result
}

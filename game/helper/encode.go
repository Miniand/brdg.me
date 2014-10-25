package helper

import (
	"bytes"
	"encoding/gob"
)

func Encode(game interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(game)
	return buf.Bytes(), err
}

func Decode(game interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(game)
}

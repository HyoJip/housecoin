package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleError(e error) {
	if e != nil {
		log.Panicln(e)
	}
}

func ToBytes(i any) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleError(encoder.Encode(i))
	return buffer.Bytes()
}

func FromBytes(i any, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	target := fmt.Sprintf("%v", i)
	hashed := sha256.Sum256([]byte(target))
	return fmt.Sprintf("%x", hashed)
}

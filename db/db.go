package db

import (
	"bytes"
	"encoding/gob"
	"github.com/peterbourgon/diskv/v3"
)

var disk *diskv.Diskv

func Init(options diskv.Options) {
	disk = diskv.New(options)
}

func Write[T any](key string, value T) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(value)
	if err != nil {
		return err
	}
	return disk.Write(key, buffer.Bytes())
}

func Read[T any](key string) (T, error) {
	var value T

	buffer, err := disk.Read(key)
	if err != nil {
		return value, err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	err = dec.Decode(&value)
	return value, err
}

package protocol

import (
	"bytes"
	"errors"
	"log"
	"math"
	"strconv"
)

type RequestType string

const (
	GET RequestType = "GET"
	SET RequestType = "SET"
	DEL RequestType = "DEL"
)

type Request struct {
	Method RequestType
	Key    string
	Value  []byte
}

type decoder struct {
	stream []byte
}

func DecryptQuery(stream []byte) (*Request, error) {
	dec := decoder{stream: stream}
	if len(stream) == 0 {
		return nil, errors.New("empty request")
	}
	cmd_, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	cmd := string(cmd_)
	if !isCmdInvalid(cmd) {
		return nil, errors.New("Invalid request")
	}
	kSize, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	if len(kSize) > math.MaxUint16 {
		return nil, errors.New("Key size should not exceed (2 Bytes)")
	}
	kSizeVal, err := strconv.Atoi(string(kSize))
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Key size is not a number")
	}
	if cmd == string(GET) || cmd == string(DEL) {
		key := string(dec.extractBytes(kSizeVal))
		return &Request{
			Method: RequestType(cmd),
			Key:    key,
			Value:  make([]byte, 0),
		}, nil
	}
	vSize, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	if len(vSize) > math.MaxUint32 {
		return nil, errors.New("Value size should not exceed (4 Bytes)")
	}
	vSizeval, err := strconv.Atoi(string(vSize))
	if err != nil {
		return nil, errors.New("Value size is not a number")
	}

	key := string(dec.extractBytes(kSizeVal))
	value := dec.extractBytes(vSizeval)

	if !isCmdInvalid(cmd) {
		return &Request{}, errors.New("Invalid command - not in (GET, SET, DEL)")
	}

	return &Request{
		Method: RequestType(cmd),
		Key:    key,
		Value:  value,
	}, nil
}

func isCmdInvalid(input string) bool {
	switch RequestType(input) {
	case GET, SET, DEL:
		return true
	default:
		return false
	}
}

func (d *decoder) extractSpaceBlock() ([]byte, error) {
	spaceIdx := bytes.IndexByte(d.stream, byte(' '))
	if spaceIdx == -1 {
		return nil, errors.New("Invalid request")
	}
	cmd := make([]byte, spaceIdx)
	copy(cmd, d.stream[:spaceIdx])
	d.stream = d.stream[spaceIdx+1:]
	return cmd, nil
}

func (d *decoder) extractBytes(size int) []byte {
	part := make([]byte, size)
	copy(part, d.stream[:size])
	d.stream = d.stream[size:]
	spaceIdx := bytes.IndexByte(d.stream, byte(' '))
	if spaceIdx == 0 {
		d.stream = d.stream[1:]
	}
	return part
}

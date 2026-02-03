package protocol

import (
	"bytes"
	"errors"
	"math"
	"strconv"
)

type RequestType string
type ProtocolError string

const (
	GET RequestType = "GET"
	SET RequestType = "SET"
	DEL RequestType = "DEL"
)

const (
	EMPTY_REQUEST ProtocolError = "empty request"
	INVALID_COMMAND ProtocolError = "Unhandled command"
	KEY_SIZE_LIMIT_EXCEEDED ProtocolError = "Key size limit exceeded"
	VALUE_SIZE_LIMIT_EXCEEDED ProtocolError = "Value size limit exceeded"
	KEY_SIZE_NOT_A_NUMBER ProtocolError = "Key size is not a number"
	VALUE_SIZE_NOT_A_NUMBER ProtocolError = "Value size is not a number"
	INVALID_KEY_SIZE ProtocolError = "Key size is invalid"
	INVALID_VALUE_SIZE ProtocolError = "Value size is invalid"
	FAILED_TO_PARSE_REQUEST ProtocolError = "Failed to parse request"
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
		return nil, errors.New(string(EMPTY_REQUEST))
	}
	rawCmd, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	cmd := string(rawCmd)
	if !isCmdInvalid(cmd) {
		return nil, errors.New(string(INVALID_COMMAND))
	}
	kSize, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	kSizeVal, err := strconv.Atoi(string(kSize))
	if kSizeVal > math.MaxUint16 {
		return nil, errors.New(string(KEY_SIZE_LIMIT_EXCEEDED))
	}
	if err != nil {
		return nil, errors.New(string(KEY_SIZE_NOT_A_NUMBER))
	}
	if cmd == string(GET) || cmd == string(DEL) {
		rawKey := dec.stream
		if len(rawKey) != kSizeVal {
			return nil, errors.New(string(INVALID_KEY_SIZE))
		}
		key := string(rawKey)
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
	vSizeVal, err := strconv.Atoi(string(vSize))
	if vSizeVal > math.MaxUint32 {
		return nil, errors.New(string(VALUE_SIZE_LIMIT_EXCEEDED))
	}
	if err != nil {
		return nil, errors.New(string(VALUE_SIZE_NOT_A_NUMBER))
	}

	rawKey, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	if len(rawKey) != kSizeVal {
		return nil, errors.New(string(INVALID_KEY_SIZE))
	}
	key := string(rawKey)

	value := dec.stream
	if len(value) != vSizeVal {
		return nil, errors.New(string(INVALID_VALUE_SIZE))
	}

	if !isCmdInvalid(cmd) {
		return &Request{}, errors.New(string(INVALID_COMMAND))
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
		return nil, errors.New(string(FAILED_TO_PARSE_REQUEST))
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

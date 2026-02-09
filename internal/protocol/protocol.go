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
	SLICE_OUT_OF_RANGE ProtocolError = "Slice out of range"
)

var ErrEmptyRequest = errors.New(string(EMPTY_REQUEST))
var ErrInvalidCommand = errors.New(string(INVALID_COMMAND))
var ErrKeySizeLimitExceeded = errors.New(string(KEY_SIZE_LIMIT_EXCEEDED))
var ErrValueSizeLimitExceeded = errors.New(string(VALUE_SIZE_LIMIT_EXCEEDED))
var ErrKeySizeNotANumber = errors.New(string(KEY_SIZE_NOT_A_NUMBER))
var ErrValueSizeNotANumber = errors.New(string(VALUE_SIZE_NOT_A_NUMBER))
var ErrInvalidKeySize = errors.New(string(INVALID_KEY_SIZE))
var ErrInvalidValueSize = errors.New(string(INVALID_VALUE_SIZE))
var ErrFailedToParseRequest = errors.New(string(FAILED_TO_PARSE_REQUEST))
var ErrSliceOutOrRange = errors.New(string(SLICE_OUT_OF_RANGE))

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
		return nil, ErrEmptyRequest
	}
	rawCmd, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	cmd := string(rawCmd)
	if !isCmdInvalid(cmd) {
		return nil, ErrInvalidCommand
	}
	kSize, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	kSizeVal, err := strconv.Atoi(string(kSize))
	if kSizeVal > math.MaxUint16 {
		return nil, ErrKeySizeLimitExceeded
	}
	if err != nil {
		return nil, ErrKeySizeNotANumber
	}
	if cmd == string(GET) || cmd == string(DEL) {
		rawKey, err := dec.extractBytes(kSizeVal)
		if err != nil {
			return nil, err
		}
		if len(rawKey) != kSizeVal {
			return nil, ErrInvalidKeySize
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
		return nil, ErrValueSizeLimitExceeded
	}
	if err != nil {
		return nil, ErrValueSizeNotANumber
	}

	rawKey, err := dec.extractSpaceBlock()
	if err != nil {
		return nil, err
	}
	if len(rawKey) != kSizeVal {
		return nil, ErrInvalidKeySize
	}
	key := string(rawKey)

	value, err := dec.extractBytes(vSizeVal)
	if err != nil {
		return nil, err
	}
	if len(value) != vSizeVal {
		return nil, ErrInvalidValueSize
	}

	if !isCmdInvalid(cmd) {
		return &Request{}, ErrInvalidCommand
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
		return nil, ErrFailedToParseRequest
	}
	res := d.stream[:spaceIdx]
	d.stream = d.stream[spaceIdx+1:]
	return res, nil
}

func (d *decoder) extractBytes(size int) ([]byte, error) {
	if len(d.stream) < size {
		return nil, ErrSliceOutOrRange
	}
	if len(d.stream) == size {
		return d.stream[0:], nil
	}
	res := d.stream[:size+1]
	d.stream = d.stream[size+1:]
	return res, nil
}

package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"testing"
)

func Test_ReturnSetRequest(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	valSize := len([]byte(value))

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	req, err := DecryptQuery(streamRawRequest)

	if err != nil {
		t.Fatalf("Unable to decrypt a client query : %s", err.Error())
	}
	if req.Method != SET || req.Key != key || !bytes.Equal(req.Value, []byte(value)) {
		t.Fatal("Invalid client query decryption")
	}
}

func Test_ReturnGetRequest(t *testing.T) {
	command := "GET"
	key := "name"

	keySize := len([]byte(key))

	rawRequest := fmt.Sprintf("%s %d %s", command, keySize, key)
	streamRawRequest := []byte(rawRequest)

	req, err := DecryptQuery(streamRawRequest)

	if err != nil {
		t.Fatalf("Unable to decrypt a client query : %s", err.Error())
	}
	if req.Method != GET || req.Key != key || !bytes.Equal(req.Value, make([]byte, 0)) {
		t.Fatal("Invalid client query decryption")
	}
}

func Test_ReturnDelRequest(t *testing.T) {
	command := "DEL"
	key := "name"

	keySize := len([]byte(key))

	rawRequest := fmt.Sprintf("%s %d %s", command, keySize, key)
	streamRawRequest := []byte(rawRequest)

	req, err := DecryptQuery(streamRawRequest)

	if err != nil {
		t.Fatalf("Unable to decrypt a client query : %s", err.Error())
	}
	if req.Method != DEL || req.Key != key || !bytes.Equal(req.Value, make([]byte, 0)) {
		t.Fatal("Invalid client query decryption")
	}
}

func Test_ReturnErrorWhenKeySizeIsIncorrect(t *testing.T) {
	command := "GET"
	key := "name"

	keySize := len([]byte(key))
	wrongKeySize := keySize - 1

	rawRequest := fmt.Sprintf("%s %d %s", command, wrongKeySize, key)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrInvalidKeySize) {
		t.Fatalf("Should return an invalid key size error : %s", err.Error())
	}
}

func Test_ReturnErrorWhenValueSizeIsIncorrect(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	valSize := len([]byte(value))
	wrongValSize := valSize - 1

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, wrongValSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrInvalidValueSize) {
		t.Fatalf("Should return an invalid value size error : %s", err.Error())
	}
}

func Test_ShouldReturnAnInvalidCommandError(t *testing.T) {
	command := "UNHANDLED_COMMAND"
	key := "name"

	keySize := len([]byte(key))

	rawRequest := fmt.Sprintf("%s %d %s", command, keySize, key)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("Should return an invalid command error")
	}
}

func Test_ShouldReturnAnKeySizeLimitExceedError(t *testing.T) {
	command := "GET"
	key := "name"

	wrongKeySize := math.MaxUint16 + 10

	rawRequest := fmt.Sprintf("%s %d %s", command, wrongKeySize, key)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrKeySizeLimitExceeded) {
		t.Fatalf("Should return an key size limit exceeded error : %s", err.Error())
	}
}

func Test_ShouldReturnAnValueSizeLimitExceedError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	wrongValSize := math.MaxUint32 + 1

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, wrongValSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrValueSizeLimitExceeded) {
		t.Fatalf("Should return an key size limit exceeded error : %s", err.Error())
	}
}

func Test_ShouldReturnAKeySizeNotANumberError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := "z"
	valSize := len([]byte(value))

	rawRequest := fmt.Sprintf("%s %s %d %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrKeySizeNotANumber) {
		t.Fatalf("Should return an key size not a number error")
	}
}

func Test_ShouldReturnAValueSizeNotANumberError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	valSize := "z"

	rawRequest := fmt.Sprintf("%s %d %s %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrValueSizeNotANumber) {
		t.Fatalf("Should return an value size not a number error")
	}
}

func Test_ShouldReturnSliceOutOfRangeError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	valSize := len([]byte(value)) + 1

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrSliceOutOrRange) {
		t.Fatalf("Should return a Slice out of range error")
	}
}

func Test_ShouldReturnAFailedToParseRequestError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len([]byte(key))
	valSize := len([]byte(value))

	rawRequest := fmt.Sprintf("%s %d %d%s%s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if !errors.Is(err, ErrFailedToParseRequest) {
		t.Fatalf("Should return an failed to parse request error")
	}
}

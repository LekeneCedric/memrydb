package protocol

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func Test_ReturnSetRequest(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len(key)
	valSize := len(value)

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

	keySize := len(key)

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

	keySize := len(key)

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

	keySize := len(key)
	wrongKeySize := keySize - 1

	rawRequest := fmt.Sprintf("%s %d %s", command, wrongKeySize, key)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(INVALID_KEY_SIZE) {
		t.Fatalf("Should return an invalid key size error")
	}
}

func Test_ReturnErrorWhenValueSizeIsIncorrect(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len(key)
	valSize := len(value)
	wrongValSize := valSize - 1

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, wrongValSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(INVALID_VALUE_SIZE) {
		t.Fatalf("Should return an invalid value size error")
	}
}

func Test_ShouldReturnAnInvalidCommandError(t *testing.T) {
	command := "UNHANDLED_COMMAND"
	key := "name"

	keySize := len(key)

	rawRequest := fmt.Sprintf("%s %d %s", command, keySize, key)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(INVALID_COMMAND) {
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

	if err.Error() != string(KEY_SIZE_LIMIT_EXCEEDED) {
		t.Fatalf("Should return an key size limit exceeded error : %s", err.Error())
	}
}

func Test_ShouldReturnAnValueSizeLimitExceedError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len(key)
	wrongValSize := math.MaxUint32 + 10

	rawRequest := fmt.Sprintf("%s %d %d %s %s", command, keySize, wrongValSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(VALUE_SIZE_LIMIT_EXCEEDED) {
		t.Fatalf("Should return an key size limit exceeded error : %s", err.Error())
	}
}

func Test_ShouldReturnAKeySizeNotANumberError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := "z"
	valSize := len(value)

	rawRequest := fmt.Sprintf("%s %s %d %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(KEY_SIZE_NOT_A_NUMBER) {
		t.Fatalf("Should return an key size not a number error")
	}
}

func Test_ShouldReturnAValueSizeNotANumberError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len(key)
	valSize := "z"

	rawRequest := fmt.Sprintf("%s %d %s %s %s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(VALUE_SIZE_NOT_A_NUMBER) {
		t.Fatalf("Should return an value size not a number error")
	}
}

func Test_ShouldReturnAFailedToParseRequestError(t *testing.T) {
	command := "SET"
	key := "name"
	value := "{\"name\":\"cedric\"}"

	keySize := len(key)
	valSize := len(value)

	rawRequest := fmt.Sprintf("%s %d %d%s%s", command, keySize, valSize, key, value)
	streamRawRequest := []byte(rawRequest)

	_, err := DecryptQuery(streamRawRequest)

	if err.Error() != string(FAILED_TO_PARSE_REQUEST) {
		t.Fatalf("Should return an failed to parse request error")
	}
}

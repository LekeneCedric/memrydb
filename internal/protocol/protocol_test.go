package protocol

import (
	"bytes"
	"fmt"
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

func TestReturnDelRequest(t *testing.T) {
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

func Test_KeyAndValueShouldNotExceedSizeLimit(t *testing.T) {}

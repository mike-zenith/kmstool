package kmstool

import (
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"strings"
	"testing"
)

var keyId = "barkey"

type MockEncryptKMS struct {
	kmsiface.KMSAPI
	EncryptFunc func(*kms.EncryptInput) (*kms.EncryptOutput, error)
}

func (self *MockEncryptKMS) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return self.EncryptFunc(input)
}

func TestEncrypt_encodesWithBase64(t *testing.T) {
	rawBlob := []byte("foo")
	expected := base64.StdEncoding.EncodeToString(rawBlob)

	m := &MockEncryptKMS{
		EncryptFunc: func(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
			output := &kms.EncryptOutput{
				CiphertextBlob: rawBlob,
				KeyId:          &keyId,
			}
			return output, nil
		},
	}
	result, _ := Encrypt(m, keyId, expected)
	if !strings.EqualFold(expected, result) {
		t.Errorf("Failed asserting that %b equals %b", expected, result)
	}
}

func TestEncrypt_returnError(t *testing.T) {
	in := "foo"
	expected := errors.New("foo")

	m := &MockEncryptKMS{
		EncryptFunc: func(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
			return nil, expected
		},
	}

	_, resultError := Encrypt(m, keyId, in)
	if resultError != expected {
		t.Errorf("Failed asserting that %s equals returned error %s", expected, resultError)
	}
}

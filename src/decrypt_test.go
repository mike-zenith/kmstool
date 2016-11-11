package kmstool

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"testing"
)

type MockDecryptKMS struct {
	kmsiface.KMSAPI
	DecryptFunc func(*kms.DecryptInput) (*kms.DecryptOutput, error)
}

func (self *MockDecryptKMS) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return self.DecryptFunc(input)
}

func TestDecrypt(t *testing.T) {
	rawBlob := []byte("foo")
	rawBlobEncoded := base64.StdEncoding.EncodeToString([]byte("foo"))

	m := &MockDecryptKMS{
		DecryptFunc: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
			output := &kms.DecryptOutput{
				Plaintext: input.CiphertextBlob,
			}
			return output, nil
		},
	}
	result, err := Decrypt(m, rawBlobEncoded)

	if err != nil {
		t.Error(err)
	}

	if result != string(rawBlob) {
		t.Errorf("Failed assering that decrypt was called. Returned %+v", result)
	}
}

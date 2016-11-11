package kmstool

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/aws/aws-sdk-go/service/kms"
	"testing"
)

type WriterMock struct {
	written string
}

func (s *WriterMock) Write(p []byte) (n int, err error) {
	s.written = string(p)
	return len(p), nil
}

func CreateTestApp() *App {
	region := "eu-central-1"
	a := NewApp()
	a.Writer = &WriterMock{}
	a.Client = kms.New(mock.Session, &aws.Config{Region: &region})
	a.Client.Handlers.Clear()

	return a
}

func TestApp_Before_content(t *testing.T) {
	args := []string{
		"kmstool",
		"--test",
		"encrypt",
		"kms:test",
	}

	a := CreateTestApp()
	a.Run(args)

	if a.Item.Content() != "kms:test" {
		t.Errorf("Item is not created properly. Item: %+v", a.Item)
	}

	if a.Client == nil {
		t.Errorf("Client not created propery. Client: %+v", a.Client)
	}
}

func TestApp_CommandEncrypt_base64(t *testing.T) {
	keyId := "k"
	ciphertexBlob := []byte("output")

	args := []string{
		"kmstool",
		"--test",
		"--region",
		"eu-central-1",
		"--cmk",
		keyId,
		"encrypt",
		"kms:test",
	}

	a := CreateTestApp()

	a.Client.Handlers.Clear()
	a.Client.Handlers.Send.PushBack(func(r *request.Request) {
		data := r.Data.(*kms.EncryptOutput)
		data.KeyId = &keyId
		data.CiphertextBlob = ciphertexBlob
	})

	err := a.Run(args)
	if err != nil {
		t.Error(err)
	}

	expected := base64.StdEncoding.EncodeToString(ciphertexBlob)
	received := a.Writer.(*WriterMock).written

	if received != expected {
		t.Errorf("Failed asserting that stdout '%s' is the base64 encoded cipherext '%s'", received, expected)
	}
}

func TestApp_CommandEncrypt_replace(t *testing.T) {
	encryptInput := "foobar: barfoo: kms:test"
	expectedOutput := "foobar: barfoo: a21zOnRlc3Q="

	args := []string{
		"kmstool",
		"--test",
		"--region",
		"eu-central-1",
		"--cmk",
		"k",
		"encrypt",
		encryptInput,
	}

	a := CreateTestApp()

	a.Client.Handlers.Send.PushBack(func(r *request.Request) {
		data := r.Data.(*kms.EncryptOutput)
		data.KeyId = &keyId
		data.CiphertextBlob = []byte("kms:test")
	})

	err := a.Run(args)
	if err != nil {
		t.Error(err)
	}

	received := a.Writer.(*WriterMock).written

	if received != expectedOutput {
		t.Errorf("Failed assering that encrypt returned replaced text. Stdout: %s", received)
	}
}

func TestApp_CommandEncrypt_customRegex(t *testing.T) {
	r := `\s+[^:]+:\s*(?P<hit>(?P<key>.+))\s*`

	encryptInput := "    barfoo: test"
	expectedOutput := "    barfoo: a21zOnRlc3Q="

	args := []string{
		"kmstool",
		"--test",
		"--region",
		"eu-central-1",
		"--cmk",
		"k",
		"--regexp",
		r,
		"encrypt",
		encryptInput,
	}

	a := CreateTestApp()

	a.Client.Handlers.Send.PushBack(func(r *request.Request) {
		data := r.Data.(*kms.EncryptOutput)
		data.KeyId = &keyId
		data.CiphertextBlob = []byte("a21zOnRlc3Q=")
	})

	err := a.Run(args)
	if err != nil {
		t.Error(err)
	}

	received := a.Writer.(*WriterMock).written

	if received != expectedOutput {
		t.Errorf("Failed assering that encrypt used custom regexp. Stdout: %s", received)
	}
}

func TestApp_CommandDecrypt(t *testing.T) {
	decryptInput := "foobar: barfoo_secret: kms:dGVzdA=="
	expectedOutput := "foobar: barfoo_secret: test"

	args := []string{
		"kmstool",
		"--test",
		"--region",
		"eu-central-1",
		"decrypt",
		decryptInput,
	}

	a := CreateTestApp()

	a.Client.Handlers.Send.PushBack(func(r *request.Request) {
		data := r.Data.(*kms.DecryptOutput)
		data.Plaintext = []byte("test")
	})

	err := a.Run(args)
	if err != nil {
		t.Error(err)
	}

	received := a.Writer.(*WriterMock).written

	if received != expectedOutput {
		t.Errorf("Failed assering that decrypt returned replaced text. Stdout: %s", received)
	}
}

func TestApp_CommandDecrypt_customRegexp(t *testing.T) {
	r := `\s+[^:]+:\s*(?P<hit>(?P<key>.+))\s*`

	decryptInput := "    barfoo: a21zOnRlc3Q="
	expectedOutput := "    barfoo: test"

	args := []string{
		"kmstool",
		"--test",
		"--region",
		"eu-central-1",
		"--regexp",
		r,
		"decrypt",
		decryptInput,
	}

	a := CreateTestApp()

	a.Client.Handlers.Send.PushBack(func(r *request.Request) {
		data := r.Data.(*kms.DecryptOutput)
		data.Plaintext = []byte("test")
	})

	err := a.Run(args)
	if err != nil {
		t.Error(err)
	}

	received := a.Writer.(*WriterMock).written

	if received != expectedOutput {
		t.Errorf("Failed assering that decrypt used custom regexp. Stdout: %s", received)
	}
}

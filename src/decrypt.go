package kmstool

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

func Decrypt(kmsApi kmsiface.KMSAPI, rawText string) (string, error) {
	b, _ := base64.StdEncoding.DecodeString(rawText)
	decryptInput := &kms.DecryptInput{
		CiphertextBlob: b,
	}

	output, err := kmsApi.Decrypt(decryptInput)

	if err != nil {
		return "", err
	}

	return string(output.Plaintext), nil
}

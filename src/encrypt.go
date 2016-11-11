package kmstool

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"regexp"
)

var base64Reg = regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$`)

func IsBase64(b []byte) bool {
	return base64Reg.Match(b)
}

func Encrypt(kmsApi kmsiface.KMSAPI, keyId string, rawText string) (string, error) {
	var res string
	encryptInput := &kms.EncryptInput{
		KeyId:     &keyId,
		Plaintext: []byte(rawText),
	}
	output, err := kmsApi.Encrypt(encryptInput)
	if err != nil {
		return "", err
	}

	if !IsBase64(output.CiphertextBlob) {
		res = base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	} else {
		res = string(output.CiphertextBlob)
	}

	return res, nil
}

package kms

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DecryptBytes will decrypt the encrypted bytes provided
func DecryptBytes(svc *kms.KMS, data []byte) ([]byte, string, error) {
	contentType := http.DetectContentType(data)
	if contentType == "application/octet-stream" {
		input := &kms.DecryptInput{
			CiphertextBlob: data,
		}
		result, err := svc.Decrypt(input)
		if err != nil {
			return nil, "", err
		}
		return result.Plaintext, *result.KeyId, nil
	}
	return data, "", nil
}

// EncryptBytes will encrypt the bytes with key
func EncryptBytes(svc *kms.KMS, plaintextData []byte, keyID string) ([]byte, error) {
	input := &kms.EncryptInput{
		Plaintext: plaintextData,
		KeyId:     aws.String(keyID),
	}
	result, err := svc.Encrypt(input)
	if err != nil {
		return nil, err
	}
	return result.CiphertextBlob, nil
}

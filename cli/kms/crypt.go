package kms

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DecryptBytes will decrypt the encrypted bytes provided
func DecryptBytes(svc *kms.KMS, encryptedData []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedData,
	}
	result, err := svc.Decrypt(input)
	if err != nil {
		return nil, err
	}
	return result.Plaintext, nil
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

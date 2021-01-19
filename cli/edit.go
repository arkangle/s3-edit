package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	myKMS "github.com/tsub/s3-edit/cli/kms"
	myS3 "github.com/tsub/s3-edit/cli/s3"
	"github.com/tsub/s3-edit/config"
)

// Edit directly a file on S3
func Edit(path myS3.Path, params *config.AWSParams, kmsID string) {
	svcS3 := s3.New(params.Session)
	svcKMS := kms.New(params.Session)

	object := myS3.GetObject(svcS3, path)
	data, dataKmsID := decryptIfKMS(svcKMS, object.Body)
	tempDirPath, tempfilePath := createTempfile(path, data)
	defer os.RemoveAll(tempDirPath)
	if kmsID == "" && dataKmsID != "" {
		kmsID = dataKmsID
	} else if kmsID == "nil" {
		kmsID = ""
	}
	editedBody := editFile(tempfilePath)
	object.Body = encryptIfKMS(svcKMS, []byte(editedBody), kmsID)
	myS3.PutObject(svcS3, path, object)
}

func decryptIfKMS(svc *kms.KMS, data []byte) ([]byte, string) {
	decrypted, kmsID, err := myKMS.DecryptBytes(svc, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return decrypted, kmsID
}

func encryptIfKMS(svc *kms.KMS, data []byte, kmsID string) []byte {
	if kmsID != "" {
		encrypted, err := myKMS.EncryptBytes(svc, data, kmsID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return encrypted
	}
	return data
}

func createTempfile(path myS3.Path, body []byte) (tempDirPath string, tempfilePath string) {
	tempDirPath, err := ioutil.TempDir("/tmp", "s3-edit")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keys := strings.Split(path.Key, "/")
	fileName := keys[len(keys)-1]
	tempfilePath = tempDirPath + "/" + fileName

	if err := ioutil.WriteFile(tempfilePath, body, os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

func editFile(path string) string {
	command := getDefaultEditor() + " " + path

	cmd := exec.Command("sh", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	changedFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(changedFile[:])
}

func getDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "vi"
	}
	return editor
}

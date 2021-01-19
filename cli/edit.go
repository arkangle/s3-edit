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
	data := object.Body
	if kmsID != "" {
		decrypted, err := myKMS.DecryptBytes(svcKMS, data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = decrypted
	}
	tempDirPath, tempfilePath := createTempfile(path, data)
	defer os.RemoveAll(tempDirPath)

	editedBody := editFile(tempfilePath)
	if kmsID != "" {
		encrypted, err := myKMS.EncryptBytes(svcKMS, []byte(editedBody), kmsID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		object.Body = encrypted
	} else {
		object.Body = []byte(editedBody)
	}
	myS3.PutObject(svcS3, path, object)
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

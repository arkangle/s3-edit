package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tsub/s3-edit/cli"
	"github.com/tsub/s3-edit/cli/s3"
	"github.com/tsub/s3-edit/config"
)

var editCmd = &cobra.Command{
	Use:   "edit [S3 file path] (Kms-ID)",
	Short: "Edit directly a file on S3",
	Long:  "Edit directly a file on S3",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := s3.ParsePath(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		params, err := config.NewAWSParams(awsProfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		kmsID := ""
		if len(args) == 2 {
			kmsID = args[1]
		}
		cli.Edit(path, params, kmsID)
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
}

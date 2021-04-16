package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yoanyombapro1234/FeelGuuds/src/services/cart_service/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: "Prints servicecli version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(version.VERSION)
		return nil
	},
}

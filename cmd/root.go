package cmd

import (
	"github.com/MasoudHeydari/eps-api/delivery"
	"github.com/spf13/cobra"
)

const configFilePath = "/etc/eps/example.json"

var RootCmd = &cobra.Command{
	Use:  "eps",
	RunE: startServer,
}

func startServer(_ *cobra.Command, _ []string) error {
	return delivery.Start(configFilePath)
}

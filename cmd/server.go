package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/tinkershack/meteomunch/server"
)

// serverCmd operates the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server serves meteo data for grpc clients",
	Long:  `server serves meteo data for grpc clients`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("server called")
		server.Serve(context.Background(), args)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

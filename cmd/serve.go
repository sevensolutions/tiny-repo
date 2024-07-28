package cmd

import (
	"github.com/sevensolutions/tiny-repo/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the add command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run TinyRepo as server",
	Long:  `Run TinyRepo as server`,
	Run: func(cmd *cobra.Command, args []string) {
		server := &server.Server{}
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

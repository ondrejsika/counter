package root

import (
	"github.com/ondrejsika/counter/server"
	"github.com/ondrejsika/counter/version"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "counter",
	Short: "counter, " + version.Version,
	Run: func(c *cobra.Command, args []string) {
		server.Server()
	},
}

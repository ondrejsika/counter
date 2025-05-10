package migrate

import (
	"github.com/ondrejsika/counter/cmd/root"
	"github.com/ondrejsika/counter/migrate"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs database migrations",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		migrate.Migrate()
	},
}

func init() {
	root.Cmd.AddCommand(Cmd)
}

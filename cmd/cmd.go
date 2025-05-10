package cmd

import (
	"github.com/ondrejsika/counter/cmd/root"
	_ "github.com/ondrejsika/counter/cmd/version"
	"github.com/spf13/cobra"
)

func Execute() {
	cobra.CheckErr(root.Cmd.Execute())
}

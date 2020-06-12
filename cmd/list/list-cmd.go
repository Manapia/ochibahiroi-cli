package list

import "github.com/spf13/cobra"

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:   "list",
		Short: "Controls custom download list",
	}

	initInitCmd()
	Cmd.AddCommand(initCmd)
}

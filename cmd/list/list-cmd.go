package list

import "github.com/spf13/cobra"

var Cmd *cobra.Command

func init() {
	Cmd = &cobra.Command{
		Use:   "list",
		Short: "Controls custom download list",
	}

	initInitCmd()
	initAppendCmd()
	initListRemoveCmd()
	initReformatCmd()
	initRunCmd()

	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(appendCmd)
	Cmd.AddCommand(removeCmd)
	Cmd.AddCommand(reformatCmd)
	Cmd.AddCommand(runCmd)
}

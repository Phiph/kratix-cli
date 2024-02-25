package cmd

import (
	"github.com/phiph/kratix-cli/pkg/promise"
	"github.com/spf13/cobra"
)

type PromiseOptions struct {
	test bool
	new  bool
}

func defaultPromiseOptions() *PromiseOptions {
	return &PromiseOptions{}
}

func newInitCmd() *cobra.Command {

	promiseCmd := &cobra.Command{
		Use:          "promise",
		Short:        "Test and Develop Kratix Promises",
		SilenceUsage: true,
	}

	promiseTestCmd := &cobra.Command{
		Use:   `test`,
		Short: "Test a promise",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runTestPromise,
	}

	promiseCmd.AddCommand(promiseTestCmd)

	return promiseCmd
}

func runTestPromise(cmd *cobra.Command, args []string) error {
	var arg string
	if len(args) > 0 {
		arg = args[0]
	}
	promise.TestPromise(arg)
	return nil
}

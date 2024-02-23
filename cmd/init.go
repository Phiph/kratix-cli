package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/golang-cli-template/internal/convert"
	"github.com/FalcoSuessgott/golang-cli-template/pkg/example"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	multiply bool
	add      bool
}

func defaultInitOptions() *InitOptions {
	return &InitOptions{}
}

func newInitCmd() *cobra.Command {
	o := defaultExampleOptions()

	cmd := &cobra.Command{
		Use:          "init",
		Short:        "initialise Kratix",
		SilenceUsage: true,
		RunE:         o.run,
	}

	cmd.Flags().BoolVarP(&o.multiply, "multiply", "m", o.multiply, "multiply")
	cmd.Flags().BoolVarP(&o.add, "add", "a", o.add, "add")

	return cmd
}

func (o *InitOptions) run(cmd *cobra.Command, args []string) error {
	values, err := o.parseArgs(args)
	if err != nil {
		return err
	}

	if o.multiply {
		fmt.Fprintf(cmd.OutOrStdout(), "%d\n", example.Multiply(values[0], values[1]))
	}

	if o.add {
		fmt.Fprintf(cmd.OutOrStdout(), "%d\n", example.Add(values[0], values[1]))
	}

	return nil
}

func (o *InitOptions) parseArgs(args []string) ([]int, error) {
	values := make([]int, 2) //nolint: gomnd

	for i, a := range args {
		v, err := convert.ToInteger(a)
		if err != nil {
			return nil, fmt.Errorf("error converting to integer: %w", err)
		}

		values[i] = v
	}

	return values, nil
}

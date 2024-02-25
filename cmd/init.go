package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/phiph/kratix-cli/internal/convert"
	"github.com/phiph/kratix-cli/pkg/example"
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
	o := defaultPromiseOptions()

	cmd := &cobra.Command{
		Use:          "promise",
		Short:        "initialise Kratix",
		SilenceUsage: true,
		RunE:         o.run,
	}

	cmd.Flags().BoolVarP(&o.test, "test", "t", o.test, "test")
	cmd.Flags().BoolVarP(&o.new, "new", "n", o.new, "new")

	return cmd
}

func (o *PromiseOptions) run(cmd *cobra.Command, args []string) error {
	values, err := o.parseArgs(args)
	if err != nil {
		return err
	}

	if o.test {

		// Run a docker command with some env vats
		// docker
		// 	--env-file .env

		//docker build \
		//--tag "${PIPELINE_NAME}" \
		//"${PWD}/configure-pipeline" ;;

		// Get the Current Directory
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Get the Pipeline Name
		var pipelineName = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1]
		fmt.Println("Pipeline Name", pipelineName)

		var buildContextDirectory = dir + "/internal/configure-pipeline"
		var pwd = dir + "/internal"
		fmt.Println("PWD:", pwd)

		buildImage(buildContextDirectory, pipelineName)
		println("Image Built")

		// Run the docker command
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		// cli.ImagePull is asynchronous.
		// The reader needs to be read completely for the pull operation to complete.
		// If stdout is not required, consider using io.Discard instead of os.Stdout.

		///

		//		docker run \
		//  -v ${PWD}/internal/configure-pipeline/test-input:/kratix/input \
		//  -v ${PWD}/internal/configure-pipeline/test-output:/kratix/output $PIPELINE_NAME
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "docker.io/library/" + pipelineName + ":latest",
			//Cmd:   []string{"echo", "hello world"},
			Tty: false,
		}, &container.HostConfig{
			Binds:      []string{pwd + "/configure-pipeline/test-input:/kratix/input", pwd + "/configure-pipeline/test-output:/kratix/output"},
			AutoRemove: true,
		}, nil, nil, pipelineName)
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			panic(err)
		}

	}

	if o.new {
		fmt.Fprintf(cmd.OutOrStdout(), "%d\n", example.Add(values[0], values[1]))
	}

	return nil
}

func buildImage(buildContextDirectory string, pipelineName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	buildContext, err := archive.TarWithOptions(buildContextDirectory, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		buildContext,
		types.ImageBuildOptions{
			Context:    buildContext,
			Dockerfile: "Dockerfile",
			Tags:       []string{pipelineName + ":latest"},
		})
	if err != nil {
		panic(err)
	}
	defer imageBuildResponse.Body.Close()
}

func (o *PromiseOptions) parseArgs(args []string) ([]int, error) {
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

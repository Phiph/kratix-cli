package promise

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func TestPromise(runCommand string) {
	// Get the Current Directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
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

	var config = &container.Config{
		Image: "configure-pipeline-" + pipelineName + ":latest",
		Tty:   false,
	}
	if runCommand != "" {
		config.Entrypoint = []string{runCommand}
	}

	//		docker run \
	//  -v ${PWD}/internal/configure-pipeline/test-input:/kratix/input \
	//  -v ${PWD}/internal/configure-pipeline/test-output:/kratix/output $PIPELINE_NAME
	resp, err := cli.ContainerCreate(ctx, config, &container.HostConfig{
		Binds: []string{
			pwd + "/configure-pipeline/test-input:/kratix/input",
			pwd + "/configure-pipeline/test-output:/kratix/output"},
		AutoRemove: true,
	}, nil, nil, pipelineName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
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

	fmt.Println(pipelineName + ":latest")

	tags := []string{"configure-pipeline-" + pipelineName + ":latest"}

	fmt.Print(tags)
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		buildContext,
		types.ImageBuildOptions{
			Context:    buildContext,
			Dockerfile: "Dockerfile",
			Tags:       tags,
		})
	if err != nil {
		panic(err)
	}
	defer imageBuildResponse.Body.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	imageExists := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == tags[0] {
				imageExists = true
				break
			}
		}
	}

	if imageExists {
		fmt.Println("Image template:latest exists.")
	} else {
		fmt.Println("Image template:latest does not exist.")
	}

}

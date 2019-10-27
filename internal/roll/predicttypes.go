package roll

import (
	"context"
	"fmt"
	"os"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type ContainerName string

func (*ContainerName) Complete(args string) []string {
	cli, err := docker.NewEnvClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCompletion error - can't get docker client - err: %v\n", err)
		return nil
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCompletion error - can't list container - err: %v\n", err)
		return nil
	}
	ret := make([]string, len(containers))
	for i, container := range containers {
		if container.Names[0][0] == '/' {
			ret[i] = container.Names[0][1:]
		} else {
			ret[i] = container.Names[0]
		}
	}
	return ret
}
func (p *ContainerName) Set(val string) (err error) {
	*p = ContainerName(val)
	return
}
func (p ContainerName) String() string {
	return string(p)
}

type ImageName string

func (*ImageName) Complete(args string) []string {
	cli, err := docker.NewEnvClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCompletion error - can't get docker client - err: %v\n", err)
		return nil
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nCompletion error - can't list container - err: %v\n", err)
		return nil
	}
	ret := make([]string, 0)
	for _, img := range images {
		ret = append(ret, img.RepoTags...)
	}
	return ret
}
func (p *ImageName) Set(val string) (err error) {
	*p = ImageName(val)
	return
}
func (p ImageName) String() string {
	return string(p)
}

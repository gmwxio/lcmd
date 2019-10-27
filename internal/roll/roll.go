package roll

import (
	"context"
	"fmt"
	"strings"

	"github.com/jpillora/opts"

	"github.com/golang/glog"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

type roll struct {
	Container ContainerName //`opts:"mode=arg"`
	Image     ImageName     //`opts:"mode=arg"`
}

func Register(opt opts.Opts) opts.Opts {
	opt.AddCommand(opts.New(&roll{}).Name("roll"))
	opt.AddCommand(opts.New(&spinup{}).Name("new"))
	return opt
}

func (rl *roll) Run() error {
	fmt.Printf("--%+v\n", rl)
	// return nil
	cli, err := docker.NewEnvClient()
	if err != nil {
		glog.Warning(err)
		return err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		glog.Warning(err)
		return err
	}
	var co *types.Container
	val := rl.Container.String()
	if strings.Index(val, "/") == -1 {
		val = "/" + val
	}
	for _, container := range containers {
		if container.Names[0] == val {
			co = &container
			break
		}
	}
	if co == nil {
		glog.Warningf("no container found with name '%s'", rl.Container)
		return fmt.Errorf("no container found with name '%s'", rl.Container)
	}
	ci, err := cli.ContainerInspect(context.Background(), co.ID)
	if err != nil {
		glog.Warning(err)
		return err
	}
	spinup := &spinup{
		Name:   rl.Container.String(),
		Image:  rl.Image,
		ports:  ci.NetworkSettings.Ports,
		env:    ci.Config.Env,
		labels: ci.Config.Labels,
		capAdd: ci.HostConfig.CapAdd,
		mounts: ci.HostConfig.Mounts,
	}
	fmt.Printf("rolling %+v\n", spinup)
	_ = cli.ContainerRemove(context.Background(), co.ID, types.ContainerRemoveOptions{Force: true})
	return spinup.Run()
}

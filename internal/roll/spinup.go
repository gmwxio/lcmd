package roll

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/network"
	volumetypes "docker.io/go-docker/api/types/volume"
	"github.com/docker/go-connections/nat"
	"github.com/golang/glog"
)

type spinup struct {
	Name         string `type:"arg" help:"name to be set on the container"`
	HostMount    bool   `help:"mount /var/run into the container"`
	NxPort       int    `help:"if not set then random via publish all"`
	Image        ImageName
	Username     string
	PasswdHash   string
	Gid          int
	Uid          int
	HomeBindPath string `help:"if set then path in mounted as /home else a volume with the same name as the container is used"`
	CapSysAdmin  bool
	// HostDockerGid int
	// internal
	ports  nat.PortMap
	env    []string
	labels map[string]string
	capAdd []string
	mounts []mount.Mount
}

func (rl *spinup) Run() error {
	if rl.Image == "" {
		return fmt.Errorf("No image specified")
	}
	cli, err := docker.NewEnvClient()
	if err != nil {
		glog.Warning(err)
		return err
	}
	config := &container.Config{
		Env:    rl.env,
		Image:  rl.Image.String(),
		Labels: rl.labels,
		// Env: []string{
		// 	"U_UID=432",
		// 	"U_GID=434",
		// 	"U_PASSWDHASH=$1$h9MkPIDe$pMhfKOfZqkqLFoNzADAMb.",
		// 	"U_NAME=garym",
		// },
		// Image:  "xfce",
		// Labels: map[string]string{"a_lable_key": "b_l_value"},
	}
	hostConfig := &container.HostConfig{
		// PortBindings: nat.PortMap{
		// 	"22/tcp":   []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "32782"}},
		// 	"4000/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "32781"}},
		// 	"4080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "32780"}},
		// 	"4443/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "32779"}},
		// 	"8080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "32778"}},
		// },
		ShmSize:       536870912,
		RestartPolicy: container.RestartPolicy{Name: "no"},
		AutoRemove:    true,
		// CapAdd:        []string{"SYS_PTRACE"},
		// Mounts: []mount.Mount{
		// 	{
		// 		Type:     "volume",
		// 		Source:   "garym",
		// 		Target:   "/home",
		// 		ReadOnly: false,
		// 		// Consistency:
		// 	},
		// 	{
		// 		Type:     "bind",
		// 		Source:   "/var/run",
		// 		Target:   "/var/run/host/",
		// 		ReadOnly: false,
		// 		// Consistency:
		// 	},
		// },
		// /var/run:/var/run/host/
		Resources: container.Resources{},
	}
	if len(rl.ports) != 0 {
		hostConfig.PortBindings = rl.ports
		if rl.NxPort != 0 {
			hostConfig.PortBindings["4000/tcp"] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(rl.NxPort)}}
		}
	} else {
		hostConfig.PublishAllPorts = true
		if rl.NxPort != 0 {
			hostConfig.PortBindings = nat.PortMap{
				"4000/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(rl.NxPort)}},
			}
		}
	}
	if rl.Username != "" {
		config.Env = append(config.Env, fmt.Sprintf("U_NAME=%s", rl.Username))
	}
	if rl.PasswdHash != "" {
		config.Env = append(config.Env, fmt.Sprintf("U_PASSWDHASH=%s", rl.PasswdHash))
	}
	if rl.Uid != 0 {
		config.Env = append(config.Env, fmt.Sprintf("U_UID=%d", rl.Uid))
	}
	if rl.Gid != 0 {
		config.Env = append(config.Env, fmt.Sprintf("U_GID=%d", rl.Gid))
	}
	//
	if len(rl.capAdd) != 0 {
		hostConfig.CapAdd = rl.capAdd
	} else {
		hostConfig.CapAdd = []string{"SYS_PTRACE"}
	}
	if rl.CapSysAdmin {
		hostConfig.CapAdd = append(hostConfig.CapAdd, "SYS_ADMIN")
	}
	if len(rl.mounts) != 0 {
		hostConfig.Mounts = rl.mounts
	} else {
		if rl.HomeBindPath != "" {
			hostConfig.Mounts = append(
				hostConfig.Mounts,
				mount.Mount{
					Type:     "bind",
					Source:   rl.HomeBindPath,
					Target:   "/home",
					ReadOnly: false,
					// Consistency:
				})
		} else {
			// ignoring errors
			_, _ = cli.VolumeCreate(context.Background(), volumetypes.VolumesCreateBody{
				Driver: "local",
				Name:   rl.Name,
			})
			hostConfig.Mounts = append(
				hostConfig.Mounts,
				mount.Mount{
					Type:     "volume",
					Source:   rl.Name,
					Target:   "/home",
					ReadOnly: false,
					// Consistency:
				})
		}
		//
		if rl.HostMount {
			hostConfig.Mounts = append(
				hostConfig.Mounts,
				mount.Mount{
					Type:     "bind",
					Source:   "/var/run",
					Target:   "/var/run/host/",
					ReadOnly: false,
					// Consistency:
				},
			)
			hostConfig.Mounts = append(
				hostConfig.Mounts,
				mount.Mount{
					Type:     "bind",
					Source:   "/var/run/docker.sock",
					Target:   "/var/run/docker.sock",
					ReadOnly: false,
					// Consistency:
				},
			)
		}
	}

	networkingConfig := &network.NetworkingConfig{}
	containerName := rl.Name
	cccb, err := cli.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, containerName)
	if err != nil {
		glog.Warning(err)
		return err
	}
	if len(cccb.Warnings) != 0 {
		fmt.Fprintf(os.Stderr, "\ndocker create warnings : %v\n", cccb.Warnings)
	}
	err = cli.ContainerStart(context.Background(), cccb.ID, types.ContainerStartOptions{})
	if err != nil {
		glog.Warning(err)
		return err
	}
	return nil
}

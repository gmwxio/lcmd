package roll

import (
	"context"
	"fmt"
	"strings"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"

	"github.com/golang/glog"
	pb "github.com/wxio/lcmd/lolaservice"
)

type lolasvc struct {
}

// New Lola Service
func New() pb.LolaServer {
	return &lolasvc{}
}

func (svc *lolasvc) SpinUp(ctx context.Context, req *pb.SpinUpRequest) (*pb.SpinUpResponse, error) {
	spinup := &spinup{
		Name:         req.Name,
		NxPort:       int(req.NxPort),
		Image:        ImageName(req.Image),
		Username:     req.Username,
		PasswdHash:   req.PasswdHash,
		Gid:          int(req.Gid),
		Uid:          int(req.Uid),
		HomeBindPath: req.HomeBindPath,
		CapSysAdmin:  req.CapSysAdmin,
	}
	fmt.Printf("rolling %+v\n", spinup)
	err := spinup.Run()
	if err != nil {
		return nil, err
	}
	cli, err := docker.NewEnvClient()
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	var co *types.Container
	val := req.Name
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
		glog.Warningf("no container found with name '%s'", req.Name)
		return nil, fmt.Errorf("no container found with name '%s'", req.Name)
	}
	ci, err := cli.ContainerInspect(context.Background(), co.ID)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}
	resp := &pb.SpinUpResponse{
		Ports: &pb.SpinUpResponse_PortMap{
			Ports: map[string]*pb.SpinUpResponse_Bindings{},
		},
	}
	for k, v := range ci.NetworkSettings.Ports {
		for _, b := range v {
			cur, ex := resp.Ports.Ports[string(k)]
			if !ex {
				cur = &pb.SpinUpResponse_Bindings{}
			}
			cur.Binding = append(cur.Binding,
				&pb.SpinUpResponse_Bindings_Binding{
					HostIp:   string(b.HostIP),
					HostPort: string(b.HostPort),
				})
			resp.Ports.Ports[string(k)] = cur
		}
	}
	return resp, nil
}

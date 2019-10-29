package grpcsvc

import (
	"fmt"
	"net"

	"github.com/jpillora/opts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/wxio/lcmd/internal/roll"
	pb "github.com/wxio/lcmd/lolaservice"
)

type lcmdGrpc struct {
	Port int
}

func Register(opt opts.Opts) opts.Opts {
	opt.AddCommand(opts.New(&lcmdGrpc{Port: 50051}).Name("grpc_server"))
	return opt
}

func (cmd *lcmdGrpc) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cmd.Port))
	if err != nil {
		return err
		// log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on localhost:%d\n", cmd.Port)
	svr := grpc.NewServer(
	// grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	// grpc.Creds(creds),
	)
	pb.RegisterLolaServer(svr, roll.New())
	reflection.Register(svr)
	err = svr.Serve(lis)
	return err
}

package services

import (
	"fmt"
	"net"

	pb "dancheg97.ru/dancheg97/ctlpkg/gen/proto/v1"
	"dancheg97.ru/dancheg97/ctlpkg/services/httpservers"
	"dancheg97.ru/dancheg97/ctlpkg/services/pacman"
	"dancheg97.ru/dancheg97/ctlpkg/src"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Params struct {
	HttpPort int
	GrpcPort int
	PkgPath  string
	YayPath  string
	WebPath  string
	RepoName string
	Packager *src.OsHelper
}

func Run(params *Params) error {
	server := grpc.NewServer(
		getUnaryMiddleware(),
		getStreamMiddleware(),
	)

	pacmanService := pacman.Handlers{
		Helper:   params.Packager,
		YayPath:  params.YayPath,
		PkgPath:  params.PkgPath,
		RepoName: params.RepoName,
	}
	pb.RegisterPacmanServiceServer(server, pacmanService)
	reflection.Register(server)

	go func() {
		err := httpservers.RunHttpWrapper(httpservers.Params{
			HttpPort:  params.HttpPort,
			GrpcPort:  params.GrpcPort,
			PkgsDir:   params.PkgPath,
			WebDir:    params.WebPath,
			PacmanSvc: pacmanService,
		})
		panic(err)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(`:%d`, params.GrpcPort))
	if err != nil {
		return err
	}

	logrus.Infof(`Grpc server started on port: %d`, params.GrpcPort)
	return server.Serve(lis)
}

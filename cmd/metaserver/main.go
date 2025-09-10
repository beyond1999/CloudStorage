package main

import (
	"flag"
	"net"
	"strings"

	gen "CloudStorage/pkg/api/gen"
	"CloudStorage/pkg/meta"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type svc struct {
	gen.UnimplementedMetaServiceServer
	m *meta.Client
}

func main() {
	endpoints := flag.String("etcd", "http://127.0.0.1:2379", "comma sep etcd endpoints")
	addr := flag.String("addr", ":9090", "listen addr")
	flag.Parse()

	m, err := meta.New(strings.Split(*endpoints, ","))
	if err != nil {
		log.Fatal().Err(err).Msg("meta client")
	}

	lis, _ := net.Listen("tcp", *addr)
	g := grpc.NewServer()
	gen.RegisterMetaServiceServer(g, &svc{m: m})
	log.Info().Str("addr", *addr).Msg("MetaServer listening")
	_ = g.Serve(lis)
}

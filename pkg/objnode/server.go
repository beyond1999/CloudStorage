package objnode

import (
	"context"
	"net"

	"CloudStorage/pkg/api/gen"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	gen.UnimplementedObjNodeServiceServer
	st *Storage
}

func NewServer(st *Storage) *Server { return &Server{st: st} }

func (s *Server) PutChunk(ctx context.Context, req *gen.PutChunkRequest) (*gen.PutChunkResponse, error) {
	if err := s.st.PutChunk(req.ChunkId, req.Index, req.Data); err != nil {
		return nil, err
	}
	return &gen.PutChunkResponse{}, nil
}

func (s *Server) GetChunk(ctx context.Context, req *gen.GetChunkRequest) (*gen.GetChunkResponse, error) {
	b, err := s.st.GetChunk(req.ChunkId, req.Index)
	if err != nil {
		return nil, err
	}
	return &gen.GetChunkResponse{Data: b}, nil
}

func Serve(addr string, st *Storage) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g := grpc.NewServer()
	gen.RegisterObjNodeServiceServer(g, NewServer(st))
	log.Info().Str("addr", addr).Msg("ObjNode listening")
	return g.Serve(lis)
}

package server

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/filedrive-team/filehelper/importer"
	"github.com/ipfs/go-cid"
	pb "github.com/ipfs/go-unixfs/pb"
)

// DagPoolService is used to implement DagPoolServer.
type DagPoolService struct {
	UnimplementedDagPoolServer
	DagPool *pool.DagPool
}

func (s *DagPoolService) Add(ctx context.Context, in *AddRequest) (*AddReply, error) {
	data, err := importer.NewDagWithData(in.Block, pb.Data_File, s.DagPool.CidBuilder)
	if err != nil {
		return &AddReply{Cid: ""}, err
	}
	if !s.DagPool.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &AddReply{Cid: ""}, err
	}
	err = s.DagPool.Add(ctx, data)
	if err != nil {
		return &AddReply{Cid: ""}, err
	}
	return &AddReply{Cid: data.Cid().String()}, nil
}
func (s *DagPoolService) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	if !s.DagPool.Iam.CheckUserPolicy(in.User.Username, in.User.Pass, userpolicy.OnlyWrite) {
		return &GetReply{Block: nil}, userpolicy.AccessDenied
	}
	cid, err := cid.Decode(in.Cid)
	if err != nil {
		return nil, err
	}
	get, err := s.DagPool.Get(ctx, cid)
	if err != nil {
		return &GetReply{Block: nil}, err
	}
	return &GetReply{Block: get.RawData()}, nil
}

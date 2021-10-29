package files

import (
	"context"

	"github.com/MicroFileServer/proto"
	serverbefore "github.com/MicroFileServer/service/serverbefore/grpc"
	"github.com/MicroFileServer/service/encoder"
	gt "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	deleteFile	gt.Handler

	proto.UnimplementedMicroFileServerServer
}

func NewGRPCServer(
	endpoins	Endpoints,
) proto.MicroFileServerServer {
	return &gRPCServer{
		deleteFile: gt.NewServer(
			endpoins.DeleteFile,
			GRPCDecodeDeleteFileReq,
			encoder.EncodeGRPCResponce,
			gt.ServerBefore(
				serverbefore.PutTokenIntoCTX,
			),
		),
	}
}

func (g *gRPCServer) DeleteFile(
	ctx context.Context, 
	req *proto.DeleteFileReq,
) (*proto.DeleteFileResp, error) {
	_, resp, err := g.deleteFile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*proto.DeleteFileResp), nil
}
package files

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MicroFileServer/pkg/contextvalue/subcontext"
	"github.com/MicroFileServer/pkg/statuscode"
	"github.com/go-kit/kit/endpoint"
	log "github.com/sirupsen/logrus"
)

type Endpoints struct {
	UploadFile		endpoint.Endpoint
	DownloadFile	endpoint.Endpoint
	DeleteFile		endpoint.Endpoint
	GetFile			endpoint.Endpoint
	GetFiles		endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		UploadFile: makeUploadFileEndpoint(s),
		DownloadFile: makeDownloadFileEndpoint(s),
		DeleteFile: makeDeleteFileEndpoint(s),
		GetFile: makeGetFileEndpoint(s),
		GetFiles: makeGetFilesEndpoint(s),
	}
}

func makeUploadFileEndpoint(
	s	Service,
) endpoint.Endpoint {
	return func(
		ctx			context.Context,
		request		interface{},
	) (responce interface{}, err error) {
		req := request.(*UploadFileReq)
		
		fileSender, err := subcontext.GetSubFromContext(ctx)
		if err != nil {
			log.Errorf("Don't find sub id")
			return nil, statuscode.WrapStatusError(
				fmt.Errorf("Failed to find sub id"),
				http.StatusInternalServerError,
			)
		}

		req.FileSender = fileSender

		f, err := s.UploadFile(
			ctx,
			req,
		)
		if err != nil {
			return nil, err
		}

		return &UploadFileResp{
			File: f,
		}, nil
	}
}

func makeDownloadFileEndpoint(
	s	Service,
) endpoint.Endpoint {
	return func(
		ctx			context.Context,
		request		interface{},
	) (responce interface{}, err error) {
		req := request.(*DownloadFileReq)
		
		return s.DownloadFile(
			ctx,
			req.ID,
		)
	}
}

func makeDeleteFileEndpoint(
	s	Service,
) endpoint.Endpoint {
	return func(
		ctx			context.Context,
		request		interface{},
	) (responce interface{}, err error) {
		req := request.(*DeleteFileReq)
		
		if err := s.DeleteFile(
			ctx,
			req.FileID,
		); err != nil {
			return nil, err
		}

		return &DeleteFileResp{

		}, nil
	}
}

func makeGetFileEndpoint(
	s	Service,
) endpoint.Endpoint {
	return func(
		ctx			context.Context,
		request		interface{},
	) (responce interface{}, err error) {
		req := request.(*GetFileReq)

		f, err := s.GetFile(
			ctx,
			req.FileID,
		)
		if err != nil {
			return nil, err
		}

		return &GetFileResp{
			File: f,
		}, nil
	}
}

func makeGetFilesEndpoint(
	s	Service,
) endpoint.Endpoint {
	return func(
		ctx			context.Context,
		request		interface{},
	) (responce interface{}, err error) {
		req := request.(*GetFilesReq)

		fs, err := s.GetFiles(
			ctx,
			req.Query,
		)
		if err != nil {
			return nil, err
		}

		return &GetFilesResp{
			Files: fs,
		}, nil
	}
}
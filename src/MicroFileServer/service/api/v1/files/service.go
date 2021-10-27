package files

import (
	"context"

	"github.com/MicroFileServer/pkg/models/file"
)

type Service interface {
	UploadFile(
		ctx		context.Context,
		req		*UploadFileReq,
	) (*file.File, error)

	DownloadFile(
		ctx		context.Context,
		FileID	string,
	) (*DownloadFileResp, error)

	GetFile(
		ctx		context.Context,
		FileID	string,
	) (*file.File, error)

	GetFiles(
		ctx		context.Context,
		Query	GetFilesQuery,
	) ([]*file.File, error)

	IsOwner(
		ctx		context.Context,
		FileID	string,
		UserID	string,
	) error

	DeleteFile(
		ctx 	context.Context,
		FileID	string,
	) error
}
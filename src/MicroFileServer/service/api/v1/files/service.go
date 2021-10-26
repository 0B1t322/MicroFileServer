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
}
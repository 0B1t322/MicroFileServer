package files

import (
	"context"

	"github.com/MicroFileServer/pkg/models/file"
)

type Repository interface{
	FileRepository
}

type FileRepository interface {
	UploadFile(
		ctx			context.Context,
		fileName	string,
		rawFile		[]byte,
		Metadata	file.Metadata,
	) (*file.File, error)

	DownloadFile(
		ctx		context.Context,
		FileID	string,
	) ([]byte, error)

	GetFile(
		ctx		context.Context,
		FileID	string,
	) (*file.File, error)
}
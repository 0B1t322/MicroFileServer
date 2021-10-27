package files

import (
	"context"

	"github.com/MicroFileServer/pkg/models/file"
)

type GetFilesBuilder interface {
	SetFileSender(userid string) GetFilesBuilder
	SetDateSort() GetFilesBuilder
	SetNameSort() GetFilesBuilder
}

type Builders interface {
	GetFilesBuilder() GetFilesBuilder
}

type FileRepository interface{
	Builders

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

	DeleteFile(
		ctx		context.Context,
		FileId	string,
	) error

	GetFiles(
		ctx		context.Context,
		builder	GetFilesBuilder,
	) ([]*file.File, error)
}
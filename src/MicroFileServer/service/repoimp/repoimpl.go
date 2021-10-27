package repoimp

import (
	"github.com/MicroFileServer/pkg/repositories"
	"github.com/MicroFileServer/service/repoimp/files"
)

type RepoImp struct {
	files.FileRepository
}

func New(
	repo	*repositories.Repositories,
) *RepoImp {
	return &RepoImp{
		FileRepository: &files.FilesMongoDBImp{
			Repo: repo.File,
		},
	}
}
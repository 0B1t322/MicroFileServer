package files

import (
	"context"

	"github.com/MicroFileServer/pkg/repositories/baserepo"
	"github.com/MicroFileServer/pkg/models/file"
	"github.com/kamva/mgm/v3"
)

type FilesRepository struct {
	*baserepo.BaseRepositoryMongoDB
}

func New() *FilesRepository {
	f := &FilesRepository{}

	return f
}

func (f *FilesRepository) SaveFunc(
	ctx		context.Context,
	value	interface{},
) error {
	model := getPointer(value)

	err := mgm.Coll(f.MGModel).CreateWithCtx(
		ctx,
		model,
	)
	if err != nil {
		return err
	}

	return nil
}

func (f *FilesRepository) UpdateFunc(
	ctx		context.Context,
	value	interface{},
) error {
	model := getPointer(value)

	err := mgm.Coll(f.MGModel).UpdateWithCtx(
		ctx,
		model,
	)
	if err != nil {
		return err
	}

	return nil
}



func getPointer(v interface{}) *file.FileMongoDB { 
	model, ok := v.(*file.FileMongoDB)
	if !ok {
		_m := v.(file.FileMongoDB)
		model = &_m
	}
	return model
}
package baserepo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/MicroFileServer/pkg/repositories/agregate"
	"github.com/MicroFileServer/pkg/repositories/deleter"
	"github.com/MicroFileServer/pkg/repositories/getter"
	"github.com/MicroFileServer/pkg/repositories/saver"
	"github.com/MicroFileServer/pkg/repositories/updater"
	"github.com/kamva/mgm/v3"
)

type BaseRepositoryMongoDB struct {
	saver.Saver
	getter.Getter
	deleter.Deleter
	updater.Updater
	agregate.Agregater

	MGModel	mgm.Model
}

type BaseSave interface {
	SaveFunc(context.Context, interface{}) error
}

type BaseUpdate interface {
	UpdateFunc(context.Context, interface{}) error
}

type Repository interface {
	BaseSave
	BaseUpdate
}
/*
NewMongoDB

Panic if model not implements mgm.Model or model is not pointer to struct

model should be pointer to struct

	example:

	type SomeModel struct {
		mgm.DefaultModel
	}

	type SomeRepo struct {
		*BaseRepositoryMongoDB
	}

	func New() *SomeRepo {
		s := &SomeRepo{}

		s.BaseRepositoryMongoDB = baserepo.NewMongoDB(&SomeModel{}, s)

		return s
	}

	func (s *SomeRepo) SaveFunc(_ context.Context, value inteface{}) {
		// some logic
	}

	func (s *SomeRepo) UpdateFunc(_ context.Context, value inteface{}) {
		// some logic
	}
*/
func NewMongoDB(
	model	interface{},
	repo	Repository,
) *BaseRepositoryMongoDB {
	base := &BaseRepositoryMongoDB{}

	mgModel, _model, err := getMgmModel(model)
	if err != nil {
		panic(err)
	}

	base.MGModel = mgModel

	base.Saver = saver.NewSaverByType(
		_model,
		mgModel,
		repo.SaveFunc,
	)

	base.Updater = updater.NewUpdater(
		_model,
		mgModel,
		repo.UpdateFunc,
	)

	base.Getter = getter.NewGetByType(
		mgModel,
	)

	base.Deleter = deleter.NewDeleteByType(
		mgModel,
	)

	base.Agregater = agregate.NewByType(
		mgModel,
	)


	return base
}

func getMgmModel(model interface{}) (mgModel mgm.Model, structValue interface{}, err error) {
	value := reflect.ValueOf(model)
	if value.Type().Kind() == reflect.Ptr {
		mgModel, ok := value.Interface().(mgm.Model)
		if ok {
			return mgModel, value.Elem().Interface(), nil
		} else {
			return nil, nil, fmt.Errorf("Not implement %s", "mgm.Model")
		}
	} else {
		return nil, nil, fmt.Errorf("Type should be %s", reflect.PtrTo(value.Type()))
	}
}
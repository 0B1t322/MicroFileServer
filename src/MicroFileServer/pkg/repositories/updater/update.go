package updater

import (
	"context"
	"reflect"

	"github.com/MicroFileServer/pkg/repositories/utils"
	"github.com/kamva/mgm/v3"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateFunc func(ctx context.Context, value interface{}) error

type Update struct {
	_type			mgm.Model
	updateFunc		UpdateFunc
	t				reflect.Type
	tPtr			reflect.Type
}

func NewUpdater(
	Type 		interface{},
	ModelType 	mgm.Model,
	fun			UpdateFunc,
) *Update {
	u := &Update{
		_type: ModelType,
	}

	t, err := utils.GetStructType(Type)
	if err != nil {
		log.WithFields(
			log.Fields{
				"package": "saver",
				"func": "NewUpdater",
				"err": err,
			},
		).Panic()
	}

	u.t = t

	u.updateFunc = utils.CreateCheckTypeFunc(
		t,
		fun,
	)

	return u
}


func (u *Update) Update(
	ctx		context.Context,
	v		interface{},
) error {
	return u.updateFunc(ctx, v)
}

 
// if fun is not nil call them
func (u *Update) UpdateFields(
	ctx		context.Context,
	filter	interface{},
	update	interface{},
	// Can be nil
	fun func (*mongo.UpdateResult) error,
	opts	...*options.UpdateOptions,
) error {
	result, err := mgm.Coll(u._type).UpdateOne(
		ctx,
		filter,
		update,
		opts...,
	)
	if err != nil {
		return err
	}

	if fun != nil {
		return fun(result)
	}

	return nil
}



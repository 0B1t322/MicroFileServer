package updater

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Updater interface {
	Update(
		context.Context,
		interface{},
	) error

	UpdateFields(
		ctx		context.Context,
		filter	interface{},
		update	interface{},
		fun func (*mongo.UpdateResult) error,
		opts	...*options.UpdateOptions,
	) error
}
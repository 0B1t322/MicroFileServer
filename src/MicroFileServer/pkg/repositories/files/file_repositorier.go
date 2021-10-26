package files

import (
	"github.com/MicroFileServer/pkg/repositories/agregate"
	"github.com/MicroFileServer/pkg/repositories/deleter"
	"github.com/MicroFileServer/pkg/repositories/getter"
	"github.com/MicroFileServer/pkg/repositories/saver"
	"github.com/MicroFileServer/pkg/repositories/updater"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepositorier interface{
	getter.Getter
	deleter.Deleter
	updater.Updater
	agregate.Agregater
	saver.Saver
	NewBucket(
		opts	...*options.BucketOptions,
	) (*gridfs.Bucket, error)
}
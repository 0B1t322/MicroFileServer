package gridfs

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GridFS struct {

}

func (g *GridFS) NewBucket(
	opts ...*options.BucketOptions,
) (*gridfs.Bucket, error) {
	_, _, db, err := mgm.DefaultConfigs()
	if err != nil {
		return nil, err
	}

	return gridfs.NewBucket(db, opts...)
}
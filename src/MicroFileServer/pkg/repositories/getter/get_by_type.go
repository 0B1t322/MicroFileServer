package getter

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"github.com/kamva/mgm/v3"
)

type GetByType struct {
	_type mgm.Model
}

func NewGetByType(
	Type mgm.Model,
) *GetByType {
	return &GetByType{
		_type: Type,
	}
}

func (g *GetByType) GetOne(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.SingleResult) error, 
	opts ...*options.FindOneOptions,
) error {
	res := mgm.Coll(g._type).FindOne(
		ctx,
		filter,
		opts...,
	)
	if res.Err() != nil {
		return res.Err()
	}
	
	if err := f(res); err != nil {
		return err
	}

	return nil
}

type GetDataWriter interface {
	WriteGetData(
		cur 		*mongo.Cursor,
		TotalResult	int64,
		opts 		...*options.FindOptions,
	)
}

// GetAllFiltered
// 
// You can get writed information by putting context that implement GetDataWriter
// 
// example:
// 
/*

	for i := 0; i < 10; i++ {
		if err := Repository.Project.Save(
			context.Background(),
			project.Project{
				CompactProject: project.CompactProject{
					Name: fmt.Sprintf("Test_Project_%v", i),
				},
				Desrcription: fmt.Sprintf("Test_Description_%v", i),
			},
		); err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
	ctx := getter.NewGetContext(context.Background())
		if err := Repository.Project.GetAllFiltered(
			ctx,
			bson.M{},
			func(c *mongo.Cursor) error {
				var result []*project.Project
				if err := c.All(
					context.Background(),
					&result,
				); err != nil {
					return err
				}
				return nil
			},
			options.Find().SetLimit(5),
		); err != nil {
			t.Log(err)
			t.FailNow()
		}

	t.Logf("%+v", ctx)

--- Output ---

&{Context:context.Background TotalResult:10 Offset:0 Limit:5 Count:5 HasMore:true}
*/ 
func (g *GetByType) GetAllFiltered(
	ctx context.Context, 
	filter interface{}, 
	f func(*mongo.Cursor) error, 
	opts ...*options.FindOptions,
) error {
	coll := mgm.Coll(g._type)
	cur, err := coll.Find(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	if dataWriter, ok := ctx.(GetDataWriter); ok {
		
		totalResult, err := coll.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		dataWriter.WriteGetData(cur, totalResult ,opts...)
	}

	if f != nil {
		if err := f(cur); err != nil {
			return err
		}
	}

	return nil
}
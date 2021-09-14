package getter

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type GetContextWriter interface {
// 	context.Context

// 	SetTotalResult(int64)
// 	SetOffset(int64)
// 	SetLimit(int64)
// 	SetCount(int64)
// 	SetHasMore(bool)
// }

// type GetContextReader interface {
// 	context.Context

// 	TotalResult() 	int64
// 	Offset()		int64
// 	Limit()			int64
// 	Count()			int64
// 	HasMore()		bool
// }

// type GetContext interface {
// 	context.Context
// 	GetContextReader
// 	GetContextWriter

// 	Reader()	GetContextReader
// }

type GetContext struct {
	context.Context

	TotalResult		int64

	Offset			int64

	Limit			int64

	Count			int64

	HasMore			bool
}

func (g *GetContext) ReadTotalResult() int64 {
	return g.TotalResult
}

func (g *GetContext) ReadOffset() int64 {
	return g.Offset
}

func (g *GetContext) ReadLimit() int64 {
	return g.Limit
}

func (g *GetContext) ReadCount() int64 {
	return g.Count
}

func (g *GetContext) ReadHasMore() bool {
	return g.HasMore
}

func NewGetContext(ctx context.Context) *GetContext {
	return &GetContext{
		Context: ctx,
	}
}

func (g *GetContext) WriteGetData(
	cur 			*mongo.Cursor,
	totalResult		int64,
	opts 			...*options.FindOptions,
) {
	op := options.MergeFindOptions(opts...)
	
	g.TotalResult = totalResult

	if op.Limit != nil {
		g.Limit = *op.Limit
	}

	if op.Skip != nil {
		g.Offset = *op.Skip
	}

	g.Count = int64(cur.RemainingBatchLength())

	if g.TotalResult - g.Offset - g.Count > 0 {
		g.HasMore = true
	}
}
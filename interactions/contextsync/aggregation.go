package contextsync

import "context"

var aggregationCtx, aggregationCancel = context.WithCancel(context.Background())

// last writer wins
func NewAggregationContext() context.Context {
	aggregationCancel()
	aggregationCtx, aggregationCancel = context.WithCancel(context.Background())
	return aggregationCtx
}

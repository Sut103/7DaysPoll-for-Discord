package poll

import "context"

var aggregationContexts = make(map[string]AggregationContext)

type AggregationContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// last writer wins
func NewAggregationContext(channelID string, messageID string) context.Context {
	if ctx, ok := aggregationContexts[channelID+messageID]; ok {
		ctx.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	aggregationContexts[channelID+messageID] = AggregationContext{ctx, cancel}
	return ctx
}

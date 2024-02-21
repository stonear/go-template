package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// merge multiple pgx.QueryTracer implementations
type QueryTracer struct {
	impl []pgx.QueryTracer
}

func NewTracer(impl ...pgx.QueryTracer) pgx.QueryTracer {
	return QueryTracer{
		impl: impl,
	}
}

func (q QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range q.impl {
		ctx = tracer.TraceQueryStart(ctx, conn, data)
	}
	return ctx
}

func (q QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range q.impl {
		tracer.TraceQueryEnd(ctx, conn, data)
	}
}

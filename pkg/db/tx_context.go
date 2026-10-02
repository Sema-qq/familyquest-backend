package db

import "context"

type txContextKey struct{}

func withTx(ctx context.Context, tx Conn) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func txFromContext(ctx context.Context) (Conn, bool) {
	tx, ok := ctx.Value(txContextKey{}).(Conn)
	return tx, ok
}

func hasTx(ctx context.Context) bool {
	_, ok := txFromContext(ctx)
	return ok
}

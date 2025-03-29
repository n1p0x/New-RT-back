package middleware

import (
	"context"
	
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const (
	_initDataKey contextKey = "init-data"
)

func withInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, _initDataKey, initData)
}

func ctxInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(_initDataKey).(initdata.InitData)
	return initData, ok
}

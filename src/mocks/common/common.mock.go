package common

import (
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
	"github.com/stretchr/testify/mock"
)

type GuardMock struct {
	mock.Mock
}

func (g *GuardMock) Use(ctx rctx.Context) {
	ctx.Next()
}

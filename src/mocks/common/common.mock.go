package common

import (
	guard "github.com/isd-sgcu/rpkm66-gateway/src/app/middleware/auth"
	"github.com/stretchr/testify/mock"
)

type GuardMock struct {
	mock.Mock
}

func (g *GuardMock) Use(ctx guard.IContext) {
	ctx.Next()
}

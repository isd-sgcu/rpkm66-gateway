package common

import (
	guard "github.com/isd-sgcu/rnkm65-gateway/src/app/middleware/auth"
	"github.com/stretchr/testify/mock"
)

type GuardMock struct {
	mock.Mock
}

func (g *GuardMock) Validate(ctx guard.IContext) {
	ctx.Next()
}

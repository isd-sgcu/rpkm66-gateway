package auth

import (
	"net/http"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/pkg/service/auth"
)

type Guard struct {
	authSvc auth.Service
}

func NewAuthGuard(authSvc auth.Service) Guard {
	return Guard{
		authSvc,
	}
}

func (g *Guard) Validate(ctx rctx.Context) bool {
	token := ctx.Token()

	if token == "" {
		return true
	}

	payload, err := g.authSvc.Validate(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ResponseInternalErr{})
		return false
	}

	ctx.StoreValue("UserId", payload.UserId)
	ctx.StoreValue("Role", payload.Role)

	return true
}

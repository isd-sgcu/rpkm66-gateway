package auth

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/auth"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/utils"
	"net/http"
)

type Guard struct {
	service  auth.IService
	excludes map[string]struct{}
}

type IContext interface {
	Token() string
	Method() string
	Path() string
	StoreValue(string, string)
	JSON(int, interface{})
	Next()
}

func NewAuthGuard(s auth.IService, e map[string]struct{}) Guard {
	return Guard{
		service:  s,
		excludes: e,
	}
}

func (m *Guard) Validate(ctx IContext) {
	method := ctx.Method()
	path := ctx.Path()

	var id int32
	ids := utils.FindIntFromStr(path)
	if len(ids) > 0 {
		id = ids[0]
	}

	path = utils.FormatPath(method, path, id)
	if utils.IsExisted(m.excludes, path) {
		ctx.Next()
		return
	}

	token := ctx.Token()
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &dto.ResponseErr{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid token",
		})
		return
	}

	userId, errRes := m.service.Validate(token)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}

	ctx.StoreValue("UserId", userId.UserId)
	ctx.Next()
}

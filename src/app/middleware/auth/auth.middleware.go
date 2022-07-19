package auth

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/auth"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/utils"
	role "github.com/isd-sgcu/rnkm65-gateway/src/constant/auth"
	"net/http"
	"strings"
)

type Guard struct {
	service    auth.IService
	excludes   map[string]struct{}
	phase      string
	isValidate bool
}

type IContext interface {
	Token() string
	Method() string
	Path() string
	StoreValue(string, string)
	JSON(int, interface{})
	Next()
}

func NewAuthGuard(s auth.IService, e map[string]struct{}, p string) Guard {
	return Guard{
		service:    s,
		excludes:   e,
		phase:      p,
		isValidate: true,
	}
}

func (m *Guard) Use(ctx IContext) {
	m.Validate(ctx)

	if !m.isValidate {
		return
	}

	m.CheckConfig(ctx)

	if !m.isValidate {
		return
	}

	ctx.Next()
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
		m.isValidate = false
		return
	}

	payload, errRes := m.service.Validate(token)
	if errRes != nil {
		ctx.JSON(errRes.StatusCode, errRes)
		m.isValidate = false
		return
	}

	ctx.StoreValue("UserId", payload.UserId)
	ctx.Next()
}

func (m *Guard) CheckConfig(ctx IContext) {
	method := ctx.Method()
	path := ctx.Path()

	//check whether there is a token in path
	//if token exist, replace token with ":token"
	pathSlice := strings.Split(path, "/")
	//paths which can have a token is "/group/token"
	if pathSlice[1] == "group" {
		if len(pathSlice) > 2 && pathSlice[2] != "members" && pathSlice[2] != "leave" {
			token := pathSlice[2]
			path = strings.Replace(path, token, ":token", 1)
		}
	}

	var id int32
	ids := utils.FindIntFromStr(path)
	if len(ids) > 0 {
		id = ids[0]
	}

	path = utils.FormatPath(method, path, id)

	phses, ok := role.MapPath2Phase[path]
	if !ok {
		ctx.Next()
		return
	}

	currentPhase := m.phase
	for _, phs := range phses {
		if phs == currentPhase {
			ctx.Next()
			return
		}
	}

	ctx.JSON(http.StatusForbidden, &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden Resource",
		Data:       nil,
	})
	m.isValidate = false
}

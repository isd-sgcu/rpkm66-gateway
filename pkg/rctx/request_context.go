package rctx

import (
	"github.com/gin-gonic/gin"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/rctx"
)

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	UserID() string
	Role() string
	ID() (string, error)
	Query(string) string
	GetFormData(string) string
	Token() string
	StoreValue(string, string)
	Param(string) string
	File(string, map[string]struct{}, int64) (*dto.DecomposedFile, error)
	Next()
}

func NewGinCtx(c *gin.Context) Context {
	return &rctx.GinCtx{Ctx: c}
}

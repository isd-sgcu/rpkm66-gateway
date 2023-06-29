package rctx

import "github.com/isd-sgcu/rpkm66-gateway/src/app/dto"

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	UserID() string
	ID() (string, error)
	Query(string) string
	GetFormData(string) string
	Token() string
	Method() string
	Path() string
	StoreValue(string, string)
	Param(string) (string, error)
	File(string, map[string]struct{}, int64) (*dto.DecomposedFile, error)
	Next()
}

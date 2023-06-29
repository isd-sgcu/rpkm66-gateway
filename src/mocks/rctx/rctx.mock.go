package rctx

import (
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/stretchr/testify/mock"
)

type ContextMock struct {
	mock.Mock
	V      interface{}
	Header map[string]string
	Status int
}

func (c *ContextMock) JSON(status int, v interface{}) {
	c.V = v
	c.Status = status
}

func (c *ContextMock) Bind(v interface{}) error {
	args := c.Called(v)

	if args.Get(0) != nil {
		switch v.(type) {
		case *dto.GroupDto:
			*v.(*dto.GroupDto) = *args.Get(0).(*dto.GroupDto)
		case *dto.SelectBaan:
			*v.(*dto.SelectBaan) = *args.Get(0).(*dto.SelectBaan)
		case *dto.UserDto:
			*v.(*dto.UserDto) = *args.Get(0).(*dto.UserDto)
		case *dto.UpdateUserDto:
			*v.(*dto.UpdateUserDto) = *args.Get(0).(*dto.UpdateUserDto)
		case *dto.CheckinVerifyRequest:
			*v.(*dto.CheckinVerifyRequest) = *args.Get(0).(*dto.CheckinVerifyRequest)
		case *dto.Verify:
			*v.(*dto.Verify) = *args.Get(0).(*dto.Verify)
		case *dto.CheckinConfirmRequest:
			*v.(*dto.CheckinConfirmRequest) = *args.Get(0).(*dto.CheckinConfirmRequest)
		case *dto.ConfirmEstampRequest:
			*v.(*dto.ConfirmEstampRequest) = *args.Get(0).(*dto.ConfirmEstampRequest)
		}
	}

	return args.Error(1)
}

func (c *ContextMock) ID() (string, error) {
	args := c.Called()

	return args.String(0), args.Error(1)
}

func (c *ContextMock) Host() string {
	args := c.Called()

	return args.String(0)
}

func (c *ContextMock) UserID() string {
	args := c.Called()
	return args.String(0)
}

func (c *ContextMock) Query(key string) string {
	args := c.Called(key)
	return args.String(0)
}

func (c *ContextMock) Token() string {
	args := c.Called()

	return args.String(0)
}

func (c *ContextMock) StoreValue(key string, val string) {
	_ = c.Called(key, val)

	c.Header[key] = val
}

func (c *ContextMock) Method() string {
	args := c.Called()

	return args.String(0)
}

func (c *ContextMock) Path() string {
	args := c.Called()

	return args.String(0)
}

func (c *ContextMock) Next() {
	_ = c.Called()

	return
}

func (c *ContextMock) File(key string, allowContent map[string]struct{}, _ int64) (res *dto.DecomposedFile, err error) {
	args := c.Called(key, allowContent)

	if args.Get(0) != nil {
		res = args.Get(0).(*dto.DecomposedFile)
	}

	return res, args.Error(1)
}

func (c *ContextMock) GetFormData(key string) string {
	args := c.Called(key)

	return args.String(0)
}

func (c *ContextMock) Param(string) (string, error) {
	args := c.Called()

	return args.String(0), args.Error(1)
}

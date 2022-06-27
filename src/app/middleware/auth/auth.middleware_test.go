package auth

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant"
	"github.com/isd-sgcu/rnkm65-gateway/src/mocks/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthGuardTest struct {
	suite.Suite
	ExcludePath     map[string]struct{}
	UserId          string
	Token           string
	UnauthorizedErr *dto.ResponseErr
	ServiceDownErr  *dto.ResponseErr
}

func TestAuthGuard(t *testing.T) {
	suite.Run(t, new(AuthGuardTest))
}

func (u *AuthGuardTest) SetupTest() {
	u.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	u.UnauthorizedErr = &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid token",
		Data:       nil,
	}

	u.Token = faker.Word()
	u.UserId = faker.UUIDDigit()

	u.ExcludePath = map[string]struct{}{
		"POST /exclude":     {},
		"POST /exclude/:id": {},
	}
}

func (u *AuthGuardTest) TestValidateSuccess() {
	want := u.UserId

	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(&dto.TokenPayloadAuth{
		UserId: u.UserId,
		Role:   constant.USER,
	}, nil)
	c.On("StoreValue", "UserId", u.UserId)
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	actual := c.Header["UserId"]

	assert.Equal(u.T(), want, actual)
	c.AssertNumberOfCalls(u.T(), "Next", 1)
}

func (u *AuthGuardTest) TestValidateSkippedFromExcludePath() {
	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/exclude")
	c.On("Token").Return("")
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	c.AssertNumberOfCalls(u.T(), "Next", 1)
	c.AssertNumberOfCalls(u.T(), "Token", 0)
}

func (u *AuthGuardTest) TestValidateSkippedFromExcludePathWithID() {
	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/exclude/1")
	c.On("Token").Return("")
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	c.AssertNumberOfCalls(u.T(), "Next", 1)
	c.AssertNumberOfCalls(u.T(), "Token", 0)
}

func (u *AuthGuardTest) TestValidateFailed() {
	want := u.UnauthorizedErr

	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(nil, u.UnauthorizedErr)

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
}

func (u *AuthGuardTest) TestValidateTokenNotIncluded() {
	want := u.UnauthorizedErr

	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return("")
	srv.On("Validate")

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
	srv.AssertNumberOfCalls(u.T(), "Validate", 0)
}

func (u *AuthGuardTest) TestValidateTokenGrpcErr() {
	want := u.ServiceDownErr

	srv := new(auth.ServiceMock)
	c := new(auth.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(nil, u.ServiceDownErr)

	h := NewAuthGuard(srv, u.ExcludePath)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
}

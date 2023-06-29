package auth

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/src/config"
	role "github.com/isd-sgcu/rpkm66-gateway/src/constant/auth"
	"github.com/isd-sgcu/rpkm66-gateway/src/mocks/auth"
	"github.com/isd-sgcu/rpkm66-gateway/src/mocks/rctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthGuardTest struct {
	suite.Suite
	conf            config.App
	ExcludePath     map[string]struct{}
	UserId          string
	Token           string
	UnauthorizedErr *dto.ResponseErr
	ServiceDownErr  *dto.ResponseErr
	ForbiddenErr    *dto.ResponseErr
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

	u.ForbiddenErr = &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden Resource",
		Data:       nil,
	}

	u.Token = faker.Word()
	u.UserId = faker.UUIDDigit()

	u.ExcludePath = map[string]struct{}{
		"POST /exclude":     {},
		"POST /exclude/:id": {},
	}

	u.conf = config.App{
		Port:        3000,
		Debug:       true,
		Phase:       "register",
		MaxFileSize: 10000000,
	}
}

func (u *AuthGuardTest) TestValidateSuccess() {
	want := u.UserId

	srv := new(auth.ServiceMock)
	c := &rctx.ContextMock{
		Header: map[string]string{},
	}

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(&dto.TokenPayloadAuth{
		UserId: u.UserId,
	}, nil)
	c.On("StoreValue", "UserId", u.UserId)
	c.On("StoreValue", "Role", role.USER)
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	actual := c.Header["UserId"]

	assert.Equal(u.T(), want, actual)
	c.AssertNumberOfCalls(u.T(), "Next", 1)
}

func (u *AuthGuardTest) TestValidateSkippedFromExcludePath() {
	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/exclude")
	c.On("Token").Return("")
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	c.AssertNumberOfCalls(u.T(), "Next", 1)
	c.AssertNumberOfCalls(u.T(), "Token", 0)
}

func (u *AuthGuardTest) TestValidateSkippedFromExcludePathWithID() {
	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/exclude/1")
	c.On("Token").Return("")
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	c.AssertNumberOfCalls(u.T(), "Next", 1)
	c.AssertNumberOfCalls(u.T(), "Token", 0)
}

func (u *AuthGuardTest) TestValidateFailed() {
	want := u.UnauthorizedErr

	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(nil, u.UnauthorizedErr)

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
}

func (u *AuthGuardTest) TestValidateTokenNotIncluded() {
	want := u.UnauthorizedErr

	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return("")
	srv.On("Validate")

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
	srv.AssertNumberOfCalls(u.T(), "Validate", 0)
}

func (u *AuthGuardTest) TestValidateTokenGrpcErr() {
	want := u.ServiceDownErr

	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return("POST")
	c.On("Path").Return("/auth")
	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(nil, u.ServiceDownErr)

	h := NewAuthGuard(srv, u.ExcludePath, u.conf)
	h.Validate(c)

	assert.Equal(u.T(), want, c.V)
}

func testConfigSuccess(t *testing.T, u *AuthGuardTest, conf config.App, mth string, pth string) {
	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return(mth)
	c.On("Path").Return(pth)
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath, conf)
	h.CheckConfig(c)

	c.AssertNumberOfCalls(t, "Next", 1)
}

func (u *AuthGuardTest) TestConfigSuccess() {
	u.conf.Phase = "register"
	testConfigSuccess(u.T(), u, u.conf, "GET", "/user")
	testConfigSuccess(u.T(), u, u.conf, "PUT", "/user")
	u.conf.Phase = "select"
	testConfigSuccess(u.T(), u, u.conf, "GET", "/group/1")
	testConfigSuccess(u.T(), u, u.conf, "DELETE", "/group/members/2")
	testConfigSuccess(u.T(), u, u.conf, "DELETE", "/group/leave")
	u.conf.Phase = "eventDay"
	testConfigSuccess(u.T(), u, u.conf, "POST", "/qr/checkin/verify")
	u.conf.Phase = "eStamp"
	testConfigSuccess(u.T(), u, u.conf, "POST", "/qr/estamp/confirm")
}

func testConfigFail(t *testing.T, u *AuthGuardTest, conf config.App, mth string, pth string) {
	want := u.ForbiddenErr

	srv := new(auth.ServiceMock)
	c := new(rctx.ContextMock)

	c.On("Method").Return(mth)
	c.On("Path").Return(pth)
	c.On("Next")

	h := NewAuthGuard(srv, u.ExcludePath, conf)
	h.CheckConfig(c)

	assert.Equal(t, want, c.V)
}

func (u *AuthGuardTest) TestConfigFail() {
	u.conf.Phase = "register"
	testConfigFail(u.T(), u, u.conf, "PUT", "/group")
	u.conf.Phase = "select"
	testConfigFail(u.T(), u, u.conf, "PUT", "/file/upload")
	testConfigFail(u.T(), u, u.conf, "GET", "/estamp/1")
	u.conf.Phase = "eventDay"
	testConfigFail(u.T(), u, u.conf, "PUT", "/group")
	u.conf.Phase = "emStamp"
	testConfigFail(u.T(), u, u.conf, "PUT", "/group")
}

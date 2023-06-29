package auth

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/auth"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthGuardTest struct {
	suite.Suite
	UserId          string
	Role            string
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
	u.Role = "ADMIN"
}

func (u *AuthGuardTest) TestValidateSuccess() {
	wantId := u.UserId
	wantRole := u.Role

	srv := new(auth.ServiceMock)
	c := &rctx.ContextMock{
		Header: map[string]string{},
	}

	c.On("Token").Return(u.Token)
	srv.On("Validate", u.Token).Return(&dto.TokenPayloadAuth{
		UserId: u.UserId,
		Role:   u.Role,
	}, nil)
	c.On("StoreValue", "UserId", u.UserId)
	c.On("StoreValue", "Role", u.Role)
	c.On("Next")

	h := NewAuthGuard(srv)
	h.Validate(c)

	actualId := c.Header["UserId"]
	actualRole := c.Header["Role"]

	assert.Equal(u.T(), wantId, actualId)
	assert.Equal(u.T(), wantRole, actualRole)
}

// other case is TBD

package auth

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	mock "github.com/isd-sgcu/rpkm66-gateway/mocks/auth"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceTest struct {
	suite.Suite
	Credential     *proto.Credential
	Payload        *dto.TokenPayloadAuth
	Unauthorized   *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

func (t *AuthServiceTest) SetupTest() {
	t.Credential = &proto.Credential{
		AccessToken:  faker.Word(),
		RefreshToken: faker.Word(),
		ExpiresIn:    3600,
	}

	t.Payload = &dto.TokenPayloadAuth{
		UserId: faker.UUIDDigit(),
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Not Found",
		Data:       nil,
	}

	t.Unauthorized = &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Data:       nil,
	}
}

func (t *AuthServiceTest) TestVerifyTicketSuccess() {
	want := t.Credential
	ticket := faker.Word()

	c := mock.ClientMock{}
	c.On("VerifyTicket", &proto.VerifyTicketRequest{Ticket: ticket}).Return(&proto.VerifyTicketResponse{Credential: t.Credential}, nil)

	srv := NewService(&c)

	actual, err := srv.VerifyTicket(ticket)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestVerifyTicketInvalid() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Data:       nil,
	}
	ticket := faker.Word()

	c := mock.ClientMock{}
	c.On("VerifyTicket", &proto.VerifyTicketRequest{Ticket: ticket}).Return(nil, status.Error(codes.Unauthenticated, "Invalid ticket"))

	srv := NewService(&c)

	actual, err := srv.VerifyTicket(ticket)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *AuthServiceTest) TestVerifyTicketGrpcErr() {
	want := t.ServiceDownErr
	ticket := faker.Word()

	c := mock.ClientMock{}
	c.On("VerifyTicket", &proto.VerifyTicketRequest{Ticket: ticket}).Return(nil, status.Error(codes.Unavailable, ""))

	srv := NewService(&c)

	actual, err := srv.VerifyTicket(ticket)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *AuthServiceTest) TestValidateSuccess() {
	want := t.Payload
	token := faker.Word()

	c := mock.ClientMock{}
	c.On("Validate", &proto.ValidateRequest{Token: token}).Return(&proto.ValidateResponse{
		UserId: t.Payload.UserId,
	}, nil)

	srv := NewService(&c)

	actual, err := srv.Validate(token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestValidateGrpcErr() {
	want := t.ServiceDownErr
	token := faker.Word()

	c := mock.ClientMock{}
	c.On("Validate", &proto.ValidateRequest{Token: token}).Return(nil, status.Error(codes.Unavailable, ""))

	srv := NewService(&c)

	actual, err := srv.Validate(token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *AuthServiceTest) TestRefreshTokenSuccess() {
	want := t.Credential
	token := faker.Word()

	c := mock.ClientMock{}
	c.On("RefreshToken", &proto.RefreshTokenRequest{RefreshToken: token}).Return(&proto.RefreshTokenResponse{Credential: t.Credential}, nil)

	srv := NewService(&c)

	actual, err := srv.RefreshToken(token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *AuthServiceTest) TestRefreshTokenInvalidToken() {
	want := t.Unauthorized
	token := faker.Word()

	c := mock.ClientMock{}
	c.On("RefreshToken", &proto.RefreshTokenRequest{RefreshToken: token}).Return(nil, status.Error(codes.Unauthenticated, "UnAuthorize"))

	srv := NewService(&c)

	actual, err := srv.RefreshToken(token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *AuthServiceTest) TestRefreshTokenGrpcErr() {
	want := t.ServiceDownErr
	token := faker.Word()

	c := mock.ClientMock{}
	c.On("RefreshToken", &proto.RefreshTokenRequest{RefreshToken: token}).Return(nil, status.Error(codes.Unavailable, ""))

	srv := NewService(&c)

	actual, err := srv.RefreshToken(token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

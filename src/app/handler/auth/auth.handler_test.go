package auth

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/auth"
	"github.com/isd-sgcu/rnkm65-gateway/src/mocks/user"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthHandlerTest struct {
	suite.Suite
	UserDto         *proto.User
	Credential      *proto.Credential
	Payload         *dto.TokenPayloadAuth
	BadRequestErr   *dto.ResponseErr
	UnauthorizedErr *dto.ResponseErr
	NotFoundErr     *dto.ResponseErr
	ServiceDownErr  *dto.ResponseErr
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {
	t.UserDto = &proto.User{
		Id:                    faker.UUIDDigit(),
		Firstname:             faker.FirstName(),
		Lastname:              faker.LastName(),
		Nickname:              faker.Name(),
		StudentID:             faker.Word(),
		Faculty:               faker.Word(),
		Year:                  faker.Word(),
		Phone:                 faker.Phonenumber(),
		LineID:                faker.Word(),
		Email:                 faker.Email(),
		AllergyFood:           faker.Word(),
		FoodRestriction:       faker.Word(),
		AllergyMedicine:       faker.Word(),
		Disease:               faker.Word(),
		VaccineCertificateUrl: faker.URL(),
		ImageUrl:              faker.URL(),
	}

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

	t.BadRequestErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid token",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Not found user",
		Data:       nil,
	}

	t.UnauthorizedErr = &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid ID",
	}
}

func (t *AuthHandlerTest) TestVerifyTicketSuccess() {
	want := t.Credential
	ticket := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("VerifyTicket", ticket).Return(t.Credential, nil)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		VerifyTicketDto: &dto.VerifyTicket{Ticket: ticket},
	}
	c.On("Bind", &dto.VerifyTicket{}).Return(&dto.VerifyTicket{Ticket: ticket}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.VerifyTicket(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestVerifyTicketInvalid() {
	want := t.UnauthorizedErr
	ticket := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("VerifyTicket", ticket).Return(nil, t.UnauthorizedErr)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		VerifyTicketDto: &dto.VerifyTicket{Ticket: ticket},
	}
	c.On("Bind", &dto.VerifyTicket{}).Return(&dto.VerifyTicket{Ticket: ticket}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.VerifyTicket(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestVerifyTicketGrpcErr() {
	want := t.ServiceDownErr
	ticket := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("VerifyTicket", ticket).Return(nil, t.ServiceDownErr)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		VerifyTicketDto: &dto.VerifyTicket{Ticket: ticket},
	}
	c.On("Bind", &dto.VerifyTicket{}).Return(&dto.VerifyTicket{Ticket: ticket}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.VerifyTicket(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestValidateSuccess() {
	want := t.UserDto

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}

	usrSrv := &user.ServiceMock{}
	usrSrv.On("FindOne", t.Payload.UserId).Return(t.UserDto, nil)

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Payload.UserId)

	h := NewHandler(srv, usrSrv, v)

	h.Validate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestValidateInvalidUser() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Invalid user",
		Data:       nil,
	}

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}

	usrSrv := &user.ServiceMock{}
	usrSrv.On("FindOne", t.Payload.UserId).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Payload.UserId)

	h := NewHandler(srv, usrSrv, v)

	h.Validate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestValidateGrpcErrUserService() {
	want := t.ServiceDownErr

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}

	usrSrv := &user.ServiceMock{}
	usrSrv.On("FindOne", t.Payload.UserId).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.Payload.UserId)

	h := NewHandler(srv, usrSrv, v)

	h.Validate(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestRedeemRefreshTokenSuccess() {
	want := t.Credential
	token := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("RefreshToken", token).Return(t.Credential, nil)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		RefreshTokenDto: &dto.RedeemNewToken{RefreshToken: token},
	}
	c.On("Bind", &dto.RedeemNewToken{}).Return(&dto.RedeemNewToken{RefreshToken: token}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.RefreshToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestRedeemRefreshTokenInvalid() {
	want := t.UnauthorizedErr
	token := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("RefreshToken", token).Return(nil, t.UnauthorizedErr)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		RefreshTokenDto: &dto.RedeemNewToken{RefreshToken: token},
	}
	c.On("Bind", &dto.RedeemNewToken{}).Return(&dto.RedeemNewToken{RefreshToken: token}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.RefreshToken(c)

	assert.Equal(t.T(), want, c.V)
}

func (t *AuthHandlerTest) TestRedeemRefreshTokenGrpc() {
	want := t.ServiceDownErr
	token := faker.Word()

	v, _ := validator.NewValidator()

	srv := &mock.ServiceMock{}
	srv.On("RefreshToken", token).Return(nil, t.ServiceDownErr)

	usrSrv := &user.ServiceMock{}

	c := &mock.ContextMock{
		RefreshTokenDto: &dto.RedeemNewToken{RefreshToken: token},
	}
	c.On("Bind", &dto.RedeemNewToken{}).Return(&dto.RedeemNewToken{RefreshToken: token}, nil)

	h := NewHandler(srv, usrSrv, v)

	h.RefreshToken(c)

	assert.Equal(t.T(), want, c.V)
}

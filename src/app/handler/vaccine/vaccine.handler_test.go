package vaccine

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/validator"
	"github.com/isd-sgcu/rpkm66-gateway/src/mocks/rctx"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/vaccine"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlerTest struct {
	suite.Suite
	hcert             string
	user              *proto.User
	NotFoundErr       *dto.ResponseErr
	GatewayTimeoutErr *dto.ResponseErr
	ServiceDownErr    *dto.ResponseErr
}

func TestVaccinceHandler(t *testing.T) {
	suite.Run(t, new(HandlerTest))
}

func (t *HandlerTest) SetupTest() {
	t.user = &proto.User{
		Id:              faker.UUIDDigit(),
		Title:           faker.Word(),
		Firstname:       faker.Word(),
		Lastname:        faker.Word(),
		Nickname:        faker.Word(),
		StudentID:       faker.Word(),
		Faculty:         faker.Word(),
		Year:            faker.Word(),
		Phone:           faker.Phonenumber(),
		LineID:          faker.Word(),
		Email:           faker.Email(),
		AllergyFood:     faker.Word(),
		FoodRestriction: faker.Word(),
		AllergyMedicine: faker.Word(),
		Disease:         faker.Word(),
		ImageUrl:        faker.URL(),
		CanSelectBaan:   false,
	}

	t.hcert = faker.Word()

	t.GatewayTimeoutErr = &dto.ResponseErr{
		StatusCode: http.StatusGatewayTimeout,
		Message:    "Connection timeout",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "User not found",
		Data:       nil,
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}
}

func (t *HandlerTest) TestVerifySuccess() {
	srv := new(mock.ServiceMock)
	srv.On("Verify", t.hcert, t.user.Id).Return(&dto.VaccineResponse{
		FirstName: t.user.Firstname,
		LastName:  t.user.Lastname,
		IsPassed:  true,
		Uid:       t.user.StudentID,
	}, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.user.Id, nil)
	c.On("Bind", &dto.Verify{}).Return(&dto.Verify{HCert: t.hcert}, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Verify(c)

	assert.Equal(t.T(), http.StatusNoContent, c.Status)
}

func (t *HandlerTest) TestVerifyNotFound() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("Verify", t.hcert, t.user.Id).Return(nil, t.NotFoundErr)

	c := new(rctx.ContextMock)
	c.On("UserID").Return(t.user.Id, nil)
	c.On("Bind", &dto.Verify{}).Return(&dto.Verify{HCert: t.hcert}, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Verify(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusNotFound, c.Status)
}

func (t *HandlerTest) TestVerifyGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("Verify", t.hcert, t.user.Id).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.user.Id, nil)
	c.On("Bind", &dto.Verify{}).Return(&dto.Verify{HCert: t.hcert}, nil)

	v, _ := validator.NewValidator()

	h := NewHandler(srv, v)
	h.Verify(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

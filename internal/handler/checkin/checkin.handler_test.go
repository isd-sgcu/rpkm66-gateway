package qr

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	cst "github.com/isd-sgcu/rpkm66-gateway/constant/checkin"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/validator"
	cmock "github.com/isd-sgcu/rpkm66-gateway/mocks/checkin"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/rctx"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CheckinHandlerTest struct {
	suite.Suite
	User           *proto.User
	BadRequestErr  *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	Token          string
	EventType      int
	CheckinType    int32
}

func TestCheckinHandler(t *testing.T) {
	suite.Run(t, new(CheckinHandlerTest))
}

func (t *CheckinHandlerTest) SetupTest() {
	t.User = &proto.User{
		Id:              faker.UUIDDigit(),
		Title:           faker.Word(),
		Firstname:       faker.FirstName(),
		Lastname:        faker.LastName(),
		Nickname:        faker.Name(),
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
		CanSelectBaan:   true,
	}

	t.BadRequestErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
		Data:       nil,
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.Token = faker.UUIDDigit()
	t.EventType = 5
	t.CheckinType = cst.CHECKIN
}

func (t *CheckinHandlerTest) TestCheckinVerifySuccess() {
	want := &proto.CheckinVerifyResponse{
		CheckinToken: t.Token,
		CheckinType:  t.CheckinType,
	}

	s := &cmock.ServiceMock{}
	s.On("CheckinVerify", t.User.Id, t.EventType).Return(&proto.CheckinVerifyResponse{
		CheckinToken: t.Token,
		CheckinType:  t.CheckinType,
	}, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinVerifyRequest{}).Return(&dto.CheckinVerifyRequest{
		EventType: t.EventType,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinVerify(c)

	assert.Equal(t.T(), http.StatusOK, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *CheckinHandlerTest) TestCheckinVerifyBadRequest() {
	s := &cmock.ServiceMock{}
	s.On("CheckinVerify", t.User.Id, t.EventType).Return(&proto.CheckinVerifyResponse{
		CheckinToken: t.Token,
		CheckinType:  t.CheckinType,
	}, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinVerifyRequest{}).Return(nil, errors.New(""))

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinVerify(c)

	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *CheckinHandlerTest) TestCheckinVerifyThrowInnerError() {
	s := &cmock.ServiceMock{}
	s.On("CheckinVerify", t.User.Id, t.EventType).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinVerifyRequest{}).Return(&dto.CheckinVerifyRequest{
		EventType: t.EventType,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinVerify(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
	assert.Equal(t.T(), t.ServiceDownErr, c.V)
}

func (t *CheckinHandlerTest) TestCheckinConfirmSuccess() {
	want := &proto.CheckinConfirmResponse{
		Success: true,
	}

	s := &cmock.ServiceMock{}
	s.On("CheckinConfirm", t.Token).Return(&proto.CheckinConfirmResponse{
		Success: true,
	}, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinConfirmRequest{}).Return(&dto.CheckinConfirmRequest{
		Token: t.Token,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinConfirm(c)

	assert.Equal(t.T(), http.StatusOK, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *CheckinHandlerTest) TestCheckinConfirmBadRequest() {
	s := &cmock.ServiceMock{}
	s.On("CheckinConfirm", t.Token).Return(&proto.CheckinConfirmResponse{
		Success: true,
	}, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinConfirmRequest{}).Return(nil, errors.New(""))

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinConfirm(c)

	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *CheckinHandlerTest) TestCheckinConfirmThrowInnerError() {
	s := &cmock.ServiceMock{}
	s.On("CheckinConfirm", t.Token).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.User.Id)
	c.On("Bind", &dto.CheckinConfirmRequest{}).Return(&dto.CheckinConfirmRequest{
		Token: t.Token,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.CheckinConfirm(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
	assert.Equal(t.T(), t.ServiceDownErr, c.V)
}

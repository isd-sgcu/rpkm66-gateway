package checkin

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	cst "github.com/isd-sgcu/rpkm66-gateway/constant/checkin"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/checkin"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckinServiceTest struct {
	suite.Suite
	User           *proto.User
	Token          string
	ServiceDownErr *dto.ResponseErr
	BadRequestErr  *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestCheckin(t *testing.T) {
	suite.Run(t, new(CheckinServiceTest))
}

func (t *CheckinServiceTest) SetupTest() {
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

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.BadRequestErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	}

	t.Token = faker.UUIDDigit()
}

func (t *CheckinServiceTest) TestCheckinVerifySuccess() {
	want := &proto.CheckinVerifyResponse{
		CheckinToken: t.Token,
		CheckinType:  cst.CHECKIN,
	}

	c := &checkin.ClientMock{}
	c.On("CheckinVerify", &proto.CheckinVerifyRequest{Id: t.User.Id, EventType: 0}).Return(&proto.CheckinVerifyResponse{
		CheckinToken: t.Token,
		CheckinType:  cst.CHECKIN,
	}, nil)

	serv := NewService(c)

	actual, err := serv.CheckinVerify(t.User.Id, 0)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *CheckinServiceTest) TestCheckinVerifyGrpcErr() {
	c := &checkin.ClientMock{}
	c.On("CheckinVerify", &proto.CheckinVerifyRequest{Id: t.User.Id, EventType: 0}).Return(nil, status.Error(codes.Unavailable, "Service is down"))

	serv := NewService(c)

	_, err := serv.CheckinVerify(t.User.Id, 0)

	assert.Equal(t.T(), t.ServiceDownErr, err)
}

func (t *CheckinServiceTest) TestCheckinVerifyInternalError() {
	c := &checkin.ClientMock{}
	c.On("CheckinVerify", &proto.CheckinVerifyRequest{Id: t.User.Id, EventType: 0}).Return(nil, status.Error(codes.Internal, "Internal Server Error"))

	serv := NewService(c)

	_, err := serv.CheckinVerify(t.User.Id, 0)

	assert.Equal(t.T(), t.InternalErr, err)
}

func (t *CheckinServiceTest) TestCheckinConfirmSuccess() {
	want := &proto.CheckinConfirmResponse{
		Success: true,
	}

	c := &checkin.ClientMock{}
	c.On("CheckinConfirm", &proto.CheckinConfirmRequest{Token: t.Token}).Return(want, nil)

	serv := NewService(c)

	actual, err := serv.CheckinConfirm(t.Token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *CheckinServiceTest) TestCheckinConfirmGrpcError() {
	c := &checkin.ClientMock{}
	c.On("CheckinConfirm", &proto.CheckinConfirmRequest{Token: t.Token}).Return(nil, status.Error(codes.Unavailable, "Service is down"))

	serv := NewService(c)

	actual, err := serv.CheckinConfirm(t.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), t.ServiceDownErr, err)
}

func (t *CheckinServiceTest) TestCheckinConfirmInvalidCheckinType() {
	c := &checkin.ClientMock{}
	c.On("CheckinConfirm", &proto.CheckinConfirmRequest{Token: t.Token}).Return(nil, status.Error(codes.InvalidArgument, "Invalid checkin type"))

	serv := NewService(c)

	actual, err := serv.CheckinConfirm(t.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid Checkin type",
		Data:       nil,
	}, err)
}

func (t *CheckinServiceTest) TestCheckinConfirmInvalidToken() {
	c := &checkin.ClientMock{}
	c.On("CheckinConfirm", &proto.CheckinConfirmRequest{Token: t.Token}).Return(nil, status.Error(codes.PermissionDenied, "Invalid token"))

	serv := NewService(c)

	actual, err := serv.CheckinConfirm(t.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Invalid Token",
		Data:       nil,
	}, err)
}

// Note: This include invalid user id
func (t *CheckinServiceTest) TestCheckinConfirmInternalError() {
	c := &checkin.ClientMock{}
	c.On("CheckinConfirm", &proto.CheckinConfirmRequest{Token: t.Token}).Return(nil, status.Error(codes.Internal, "Internal Server Error"))

	serv := NewService(c)

	actual, err := serv.CheckinConfirm(t.Token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), t.InternalErr, err)
}

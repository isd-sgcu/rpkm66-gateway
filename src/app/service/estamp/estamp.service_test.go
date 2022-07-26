package estamp

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/estamp"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EStampServiceTest struct {
	suite.Suite
	UId            string
	Event1         *proto.Event
	Event2         *proto.Event
	Event3         *proto.Event
	User           *proto.User
	EventType      string
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	InternalErr    *dto.ResponseErr
}

func TestEStampService(t *testing.T) {
	suite.Run(t, new(EStampServiceTest))
}

func (t *EStampServiceTest) SetupTest() {
	t.UId = faker.UUIDDigit()

	t.User = &proto.User{
		Id:              faker.UUIDDigit(),
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
		CanSelectBaan:   true,
	}

	t.Event1 = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Event2 = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Event3 = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Not found",
		Data:       nil,
	}

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	}

	t.EventType = "estamp"
}

func (t *EStampServiceTest) TestFindByIdSuccess() {
	want := &proto.FindEventByIDResponse{
		Event: t.Event1,
	}

	c := &mock.ClientMock{}

	c.On("FindEventByID", &proto.FindEventByIDRequest{
		Id: t.Event1.Id,
	}).Return(want, nil)

	serv := NewService(c)
	res, err := serv.FindEventByID(t.Event1.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), res, want)
}

func (t *EStampServiceTest) TestFindByIdUnavailable() {
	c := &mock.ClientMock{}

	c.On("FindEventByID", &proto.FindEventByIDRequest{
		Id: t.Event1.Id,
	}).Return(nil, status.Error(codes.Unavailable, "Service is down"))

	serv := NewService(c)
	res, err := serv.FindEventByID(t.Event1.Id)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.ServiceDownErr, err)
}

func (t *EStampServiceTest) TestFindByIdNotFound() {
	c := &mock.ClientMock{}

	c.On("FindEventByID", &proto.FindEventByIDRequest{
		Id: t.Event1.Id,
	}).Return(nil, status.Error(codes.NotFound, "Not found"))

	serv := NewService(c)
	res, err := serv.FindEventByID(t.Event1.Id)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.NotFoundErr, err)
}

func (t *EStampServiceTest) TestFindByIdInternal() {
	c := &mock.ClientMock{}

	c.On("FindEventByID", &proto.FindEventByIDRequest{
		Id: t.Event1.Id,
	}).Return(nil, status.Error(codes.Internal, "Internal Server Error"))

	serv := NewService(c)
	res, err := serv.FindEventByID(t.Event1.Id)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.InternalErr, err)
}

func (t *EStampServiceTest) TestFindAllEventWithTypeSuccess() {
	want := &proto.FindAllEventWithTypeResponse{
		Event: []*proto.Event{
			t.Event1,
			t.Event2,
		},
	}

	c := &mock.ClientMock{}

	c.On("FindAllEventWithType", &proto.FindAllEventWithTypeRequest{
		EventType: t.EventType,
	}).Return(want, nil)

	serv := NewService(c)
	res, err := serv.FindAllEventWithType(t.EventType)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), res, want)
}

func (t *EStampServiceTest) TestFindAllEventWithTypeUnavailable() {
	c := &mock.ClientMock{}

	c.On("FindAllEventWithType", &proto.FindAllEventWithTypeRequest{
		EventType: t.EventType,
	}).Return(nil, status.Error(codes.Unavailable, "Service is down"))

	serv := NewService(c)
	res, err := serv.FindAllEventWithType(t.EventType)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.ServiceDownErr, err)
}

func (t *EStampServiceTest) TestFindAllEventWithTypeNotFound() {
	c := &mock.ClientMock{}

	c.On("FindAllEventWithType", &proto.FindAllEventWithTypeRequest{
		EventType: t.EventType,
	}).Return(nil, status.Error(codes.NotFound, "Not found"))

	serv := NewService(c)
	res, err := serv.FindAllEventWithType(t.EventType)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.NotFoundErr, err)
}

func (t *EStampServiceTest) TestFindAllEventWithTypeInternal() {
	c := &mock.ClientMock{}

	c.On("FindAllEventWithType", &proto.FindAllEventWithTypeRequest{
		EventType: t.EventType,
	}).Return(nil, status.Error(codes.Internal, "Internal Server Error"))

	serv := NewService(c)
	res, err := serv.FindAllEventWithType(t.EventType)

	assert.Nil(t.T(), res)
	assert.Equal(t.T(), t.InternalErr, err)
}

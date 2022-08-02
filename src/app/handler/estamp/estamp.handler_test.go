package estamp

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/validator"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/estamp"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EstampHandlerTest struct {
	suite.Suite
	UId            string
	Events         []*proto.Event
	EventType      string
	BadRequestErr  *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	InternalErr    *dto.ResponseErr
	ForbiddenErr   *dto.ResponseErr
}

func TestEstampHandler(t *testing.T) {
	suite.Run(t, new(EstampHandlerTest))
}

func (t *EstampHandlerTest) SetupTest() {
	t.UId = faker.UUIDDigit()

	t.Events = make([]*proto.Event, 3)

	t.Events[0] = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Events[1] = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.Events[2] = &proto.Event{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Word(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Word(),
		Code:          faker.Word(),
	}

	t.EventType = "estamp"

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

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Not found",
		Data:       nil,
	}

	t.InternalErr = &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		Data:       nil,
	}

	t.ForbiddenErr = &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden resource",
		Data:       nil,
	}
}

func (t *EstampHandlerTest) TestFindByIdSuccess() {
	want := &proto.FindEventByIDResponse{
		Event: t.Events[0],
	}

	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(want, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return(t.Events[0].Id, nil)

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusOK, cm.Status)
	assert.Equal(t.T(), want, cm.V)
}

func (t *EstampHandlerTest) TestFindByIdBadRequest() {
	s := &mock.ServiceMock{}
	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return("", errors.New(""))

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusBadRequest, cm.Status)
}

func (t *EstampHandlerTest) TestFindByIdForbidden() {
	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(nil, t.ForbiddenErr)
	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return(t.Events[0].Id, nil)

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusForbidden, cm.Status)
}

func (t *EstampHandlerTest) TestFindByIdNotFound() {
	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(nil, t.NotFoundErr)
	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return(t.Events[0].Id, nil)

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusNotFound, cm.Status)
}

func (t *EstampHandlerTest) TestFindByIdInternal() {
	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(nil, t.InternalErr)
	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return(t.Events[0].Id, nil)

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusInternalServerError, cm.Status)
}

func (t *EstampHandlerTest) TestFindByIdUnavailable() {
	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(nil, t.ServiceDownErr)
	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	cm := &mock.ContextMock{}
	cm.On("ID").Return(t.Events[0].Id, nil)

	hdr.FindEventByID(cm)

	assert.Equal(t.T(), http.StatusServiceUnavailable, cm.Status)
}

func (t *EstampHandlerTest) TestVerifyEstampSuccess() {
	want := &proto.FindEventByIDResponse{
		Event: t.Events[0],
	}

	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(&proto.FindEventByIDResponse{
		Event: t.Events[0],
	}, nil)

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.UId)
	c.On("Bind", &dto.VerifyEstampRequest{}).Return(&dto.VerifyEstampRequest{
		EventId: t.Events[0].Id,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.VerifyEstamp(c)

	assert.Equal(t.T(), http.StatusOK, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *EstampHandlerTest) TestVerifyEstampBadRequest() {
	s := &mock.ServiceMock{}

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.UId)
	c.On("Bind", &dto.VerifyEstampRequest{}).Return(nil, errors.New(""))

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.VerifyEstamp(c)

	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *EstampHandlerTest) TestVerifyEstampInnerError() {
	s := &mock.ServiceMock{}
	s.On("FindEventByID", t.Events[0].Id).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("UserID").Return(t.UId)
	c.On("Bind", &dto.VerifyEstampRequest{}).Return(&dto.VerifyEstampRequest{
		EventId: t.Events[0].Id,
	}, nil)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.VerifyEstamp(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *EstampHandlerTest) TestFindAllEventWithTypeSuccess() {
	want := &proto.FindAllEventWithTypeResponse{
		Event: []*proto.Event{
			t.Events[0],
			t.Events[1],
		},
	}

	s := &mock.ServiceMock{}
	s.On("FindAllEventWithType", t.EventType).Return(want, nil)

	c := &mock.ContextMock{}
	c.On("Query", "eventType").Return(t.EventType)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.FindAllEventWithType(c)

	assert.Equal(t.T(), http.StatusOK, c.Status)
	assert.Equal(t.T(), want, c.V)
}

func (t *EstampHandlerTest) TestFindAllEventWithTypeInnerError() {
	s := &mock.ServiceMock{}
	s.On("FindAllEventWithType", t.EventType).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("Query", "eventType").Return(t.EventType)

	v, _ := validator.NewValidator()

	hdr := NewHandler(s, v)

	hdr.FindAllEventWithType(c)

	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

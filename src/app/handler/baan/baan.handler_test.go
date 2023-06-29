package baan

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/dto"
	mock "github.com/isd-sgcu/rpkm66-gateway/src/mocks/baan"
	"github.com/isd-sgcu/rpkm66-gateway/src/mocks/rctx"
	mockUsr "github.com/isd-sgcu/rpkm66-gateway/src/mocks/user"
	"github.com/isd-sgcu/rpkm66-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BaanHandlerTest struct {
	suite.Suite
	Baan           *proto.Baan
	User           *proto.User
	userId         string
	BindErr        *dto.ResponseErr
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestBaanHandler(t *testing.T) {
	suite.Run(t, new(BaanHandlerTest))
}

func (t *BaanHandlerTest) SetupTest() {
	t.Baan = &proto.Baan{
		Id:            faker.UUIDDigit(),
		NameTH:        faker.Word(),
		DescriptionTH: faker.Paragraph(),
		NameEN:        faker.Word(),
		DescriptionEN: faker.Paragraph(),
		Size:          proto.BaanSize_M,
		Facebook:      faker.URL(),
		FacebookUrl:   faker.URL(),
		Instagram:     faker.URL(),
		InstagramUrl:  faker.URL(),
		Line:          faker.URL(),
		LineUrl:       faker.URL(),
		ImageUrl:      faker.URL(),
	}

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
		BaanId:          t.Baan.Id,
	}

	t.userId = faker.UUIDDigit()

	t.ServiceDownErr = &dto.ResponseErr{
		StatusCode: http.StatusServiceUnavailable,
		Message:    "Service is down",
		Data:       nil,
	}

	t.NotFoundErr = &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Baan not found",
		Data:       nil,
	}

	t.BindErr = &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
	}
}

func createBaans(baan *proto.Baan) []*proto.Baan {
	var baans []*proto.Baan
	baans = append(baans, baan)

	for i := 0; i < 2; i++ {
		b := proto.Baan{
			Id:            faker.UUIDDigit(),
			NameTH:        faker.Word(),
			DescriptionTH: faker.Paragraph(),
			NameEN:        faker.Word(),
			DescriptionEN: faker.Paragraph(),
			Size:          proto.BaanSize_M,
			Facebook:      faker.URL(),
			FacebookUrl:   faker.URL(),
			Instagram:     faker.URL(),
			InstagramUrl:  faker.URL(),
			Line:          faker.URL(),
			LineUrl:       faker.URL(),
			ImageUrl:      faker.URL(),
		}

		baans = append(baans, &b)
	}

	return baans
}

func (t *BaanHandlerTest) TestFindAllBaanSuccess() {
	want := createBaans(t.Baan)

	srv := new(mock.ServiceMock)
	srv.On("FindAll").Return(want, nil)

	usrSrv := new(mockUsr.ServiceMock)

	c := &rctx.ContextMock{}

	h := NewHandler(srv, usrSrv)
	h.FindAll(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *BaanHandlerTest) TestFindOneBaan() {
	want := t.Baan

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(want, nil)

	usrSrv := new(mockUsr.ServiceMock)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv, usrSrv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *BaanHandlerTest) TestFindOneFoundErr() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, t.NotFoundErr)

	usrSrv := new(mockUsr.ServiceMock)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv, usrSrv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusNotFound, c.Status)
}

func (t *BaanHandlerTest) TestFindOneBadReqeust() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid ID",
		Data:       nil,
	}

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, want)

	usrSrv := new(mockUsr.ServiceMock)

	c := &rctx.ContextMock{}
	c.On("ID").Return("", errors.New("Cannot parse id"))

	h := NewHandler(srv, usrSrv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *BaanHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, t.ServiceDownErr)

	usrSrv := new(mockUsr.ServiceMock)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv, usrSrv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

func (t *BaanHandlerTest) TestGetUserBaanSuccess() {
	want := t.Baan

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(want, nil)

	usrSrv := new(mockUsr.ServiceMock)
	usrSrv.On("FindOne", t.userId).Return(t.User, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.userId)

	h := NewHandler(srv, usrSrv)
	h.GetUserBaan(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *BaanHandlerTest) TestGetUserBaanNotHaveBaan() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Baan not found",
	}
	t.User.BaanId = ""

	srv := new(mock.ServiceMock)
	srv.On("FindOne", "").Return(nil, want)

	usrSrv := new(mockUsr.ServiceMock)
	usrSrv.On("FindOne", t.userId).Return(t.User, nil)

	c := &rctx.ContextMock{}
	c.On("UserID").Return(t.userId)

	h := NewHandler(srv, usrSrv)
	h.GetUserBaan(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusNotFound, c.Status)
}

func (t *BaanHandlerTest) TestGetUserBaanGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, t.ServiceDownErr)

	usrSrv := new(mockUsr.ServiceMock)
	usrSrv.On("FindOne", t.userId).Return(nil, t.ServiceDownErr)

	c := &rctx.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv, usrSrv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

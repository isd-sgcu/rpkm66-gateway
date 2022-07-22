package baan

import (
	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	mock "github.com/isd-sgcu/rnkm65-gateway/src/mocks/baan"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type BaanHandlerTest struct {
	suite.Suite
	Baan           *proto.Baan
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

	c := &mock.ContextMock{}

	h := NewHandler(srv)
	h.FindAll(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *BaanHandlerTest) TestFindOneBaan() {
	want := t.Baan

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(want, nil)

	c := &mock.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusOK, c.Status)
}

func (t *BaanHandlerTest) TestFindOneFoundErr() {
	want := t.NotFoundErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, t.NotFoundErr)

	c := &mock.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv)
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

	c := &mock.ContextMock{}
	c.On("ID").Return("", errors.New("Cannot parse id"))

	h := NewHandler(srv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusBadRequest, c.Status)
}

func (t *BaanHandlerTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	srv := new(mock.ServiceMock)
	srv.On("FindOne", t.Baan.Id).Return(nil, t.ServiceDownErr)

	c := &mock.ContextMock{}
	c.On("ID").Return(t.Baan.Id, nil)

	h := NewHandler(srv)
	h.FindOne(c)

	assert.Equal(t.T(), want, c.V)
	assert.Equal(t.T(), http.StatusServiceUnavailable, c.Status)
}

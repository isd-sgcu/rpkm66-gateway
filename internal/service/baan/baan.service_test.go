package baan

import (
	"net/http"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/mocks/baan"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaanServiceTest struct {
	suite.Suite
	Baan           *proto.Baan
	NotFoundErr    *dto.ResponseErr
	ServiceDownErr *dto.ResponseErr
}

func TestBaanService(t *testing.T) {
	suite.Run(t, new(BaanServiceTest))
}

func (t *BaanServiceTest) SetupTest() {
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

func (t *BaanServiceTest) TestFindAllSuccess() {
	want := createBaans(t.Baan)

	c := &baan.ClientMock{}
	c.On("FindAllBaan", &proto.FindAllBaanRequest{}).Return(&proto.FindAllBaanResponse{Baans: want}, nil)
	srv := NewService(c)

	actual, err := srv.FindAll()

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *BaanServiceTest) TestFindOneSuccess() {
	want := t.Baan

	c := &baan.ClientMock{}
	c.On("FindOneBaan", &proto.FindOneBaanRequest{Id: t.Baan.Id}).Return(&proto.FindOneBaanResponse{Baan: want}, nil)
	srv := NewService(c)

	actual, err := srv.FindOne(t.Baan.Id)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *BaanServiceTest) TestFindOneNotFound() {
	want := t.NotFoundErr

	c := &baan.ClientMock{}
	c.On("FindOneBaan", &proto.FindOneBaanRequest{Id: t.Baan.Id}).Return(nil, status.Error(codes.NotFound, "Baan not found"))

	srv := NewService(c)

	actual, err := srv.FindOne(t.Baan.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

func (t *BaanServiceTest) TestFindOneGrpcErr() {
	want := t.ServiceDownErr

	c := &baan.ClientMock{}
	c.On("FindOneBaan", &proto.FindOneBaanRequest{Id: t.Baan.Id}).Return(nil, errors.New("Service is down"))
	srv := NewService(c)

	actual, err := srv.FindOne(t.Baan.Id)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), want, err)
}

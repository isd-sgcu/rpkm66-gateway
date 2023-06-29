package baan

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/baan"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type Service interface {
	FindAll() ([]*proto.Baan, *dto.ResponseErr)
	FindOne(string) (*proto.Baan, *dto.ResponseErr)
}

func NewService(client proto.BaanServiceClient) Service {
	return baan.NewService(client)
}

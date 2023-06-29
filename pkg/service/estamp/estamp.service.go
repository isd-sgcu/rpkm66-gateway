package estamp

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/estamp"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type Service interface {
	FindEventByID(string) (*proto.FindEventByIDResponse, *dto.ResponseErr)
	FindAllEventWithType(string) (*proto.FindAllEventWithTypeResponse, *dto.ResponseErr)
}

func NewService(client proto.EventServiceClient) Service {
	return estamp.NewService(client)
}

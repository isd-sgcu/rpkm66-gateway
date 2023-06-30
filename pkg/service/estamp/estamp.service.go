package estamp

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/estamp"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/event/v1"
)

type Service interface {
	FindEventByID(string) (*proto.FindEventByIDResponse, *dto.ResponseErr)
	FindAllEventWithType(string) (*proto.FindAllEventWithTypeResponse, *dto.ResponseErr)
}

func NewService(client proto.EventServiceClient) Service {
	return estamp.NewService(client)
}

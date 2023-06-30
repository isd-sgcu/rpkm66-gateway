package estamp

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/event/v1"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.EventServiceClient
}

func NewService(client proto.EventServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) FindEventByID(id string) (*proto.FindEventByIDResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindEventByID(ctx, &proto.FindEventByIDRequest{
		Id: id,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "estamp").
			Str("method", "find event by id").
			Msg("Error while find event by id")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res, nil
}

func (s *serviceImpl) FindAllEventWithType(eventType string) (*proto.FindAllEventWithTypeResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindAllEventWithType(ctx, &proto.FindAllEventWithTypeRequest{
		EventType: eventType,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "estamp").
			Str("method", "find event with type").
			Msg("Error while find event with type")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res, nil
}

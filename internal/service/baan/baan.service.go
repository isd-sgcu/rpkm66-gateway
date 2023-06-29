package baan

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.BaanServiceClient
}

func NewService(client proto.BaanServiceClient) *serviceImpl {
	return &serviceImpl{client: client}
}

func (s *serviceImpl) FindAll() ([]*proto.Baan, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindAllBaan(ctx, &proto.FindAllBaanRequest{})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "baan").
			Str("method", "find all baan").
			Msg("Error while find all baan")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Baans, nil
}

func (s *serviceImpl) FindOne(id string) (*proto.Baan, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindOneBaan(ctx, &proto.FindOneBaanRequest{Id: id})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "baan").
			Str("method", "find one baan").
			Msg("Error while find one baan")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Baan, nil
}

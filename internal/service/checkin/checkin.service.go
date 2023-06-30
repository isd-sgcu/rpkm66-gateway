package checkin

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/checkin/v1"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.CheckinServiceClient
}

func NewService(client proto.CheckinServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) CheckinVerify(userid string, eventType int) (*proto.CheckinVerifyResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.CheckinVerify(ctx, &proto.CheckinVerifyRequest{
		Id:        userid,
		EventType: int32(eventType),
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "checkin").
			Str("method", "checkin verify").
			Msg("Error while checkin verify")

		return nil, utils.ServiceErrorHandler(err)
	}

	return &proto.CheckinVerifyResponse{
		CheckinToken: res.CheckinToken,
		CheckinType:  res.CheckinType,
	}, nil
}

func (s *serviceImpl) CheckinConfirm(token string) (*proto.CheckinConfirmResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.CheckinConfirm(ctx, &proto.CheckinConfirmRequest{
		Token: token,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "checkin").
			Str("method", "checkin confirm").
			Msg("Error while checkin confirm")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res, nil
}

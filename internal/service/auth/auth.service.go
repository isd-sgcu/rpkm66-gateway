package auth

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/auth/auth/v1"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.AuthServiceClient
}

func NewService(client proto.AuthServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) VerifyTicket(ticket string) (*proto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.VerifyTicket(ctx, &proto.VerifyTicketRequest{Ticket: ticket})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("method", "verify ticket").
			Msg("Error while verify ticket")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Credential, nil
}

func (s *serviceImpl) Validate(token string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Validate(ctx, &proto.ValidateRequest{Token: token})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("method", "validate").
			Msg("Error while validate")

		return nil, utils.ServiceErrorHandler(err)
	}

	return &dto.TokenPayloadAuth{
		UserId: res.UserId,
		Role:   res.Role,
	}, nil
}

func (s *serviceImpl) RefreshToken(token string) (*proto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.RefreshToken(ctx, &proto.RefreshTokenRequest{RefreshToken: token})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "auth").
			Str("method", "refresh token").
			Msg("Error while refresh token")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Credential, nil
}

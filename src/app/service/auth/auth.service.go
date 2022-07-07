package auth

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type Service struct {
	client proto.AuthServiceClient
}

func NewService(client proto.AuthServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) VerifyTicket(ticket string) (*proto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.VerifyTicket(ctx, &proto.VerifyTicketRequest{Ticket: ticket})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusUnauthorized,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(err).
					Str("service", "auth").
					Str("module", "verify ticket").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "verify ticket").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{

			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Credential, nil
}

func (s *Service) Validate(token string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Validate(ctx, &proto.ValidateRequest{Token: token})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusUnauthorized,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(err).
					Str("service", "auth").
					Str("module", "validate").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "validate").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return &dto.TokenPayloadAuth{
		UserId: res.UserId,
	}, nil
}

func (s *Service) RefreshToken(token string) (*proto.Credential, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.RefreshToken(ctx, &proto.RefreshTokenRequest{RefreshToken: token})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusUnauthorized,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(err).
					Str("service", "auth").
					Str("module", "refresh token").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(err).
			Str("service", "auth").
			Str("module", "refresh token").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Credential, nil
}

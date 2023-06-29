package estamp

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	client proto.EventServiceClient
}

func NewService(client proto.EventServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindEventByID(id string) (*proto.FindEventByIDResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindEventByID(ctx, &proto.FindEventByIDRequest{
		Id: id,
	})

	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
			log.Error().
				Err(err).
				Str("service", "estamp").
				Str("module", "find_by_id").
				Msg("\"Error parsing\" error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}

		switch st.Code() {
		case codes.Unavailable:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_by_id").
				Msg("Service is down")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    "Not found",
				Data:       nil,
			}
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    "Forbidden resource",
				Data:       nil,
			}
		default:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_by_id").
				Msg("Unhandled error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return res, nil
}

func (s *Service) FindAllEventWithType(eventType string) (*proto.FindAllEventWithTypeResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindAllEventWithType(ctx, &proto.FindAllEventWithTypeRequest{
		EventType: eventType,
	})

	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
			log.Error().
				Err(err).
				Str("service", "estamp").
				Str("module", "find_all_event_with_type").
				Msg("\"Error parsing\" error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}

		switch st.Code() {
		case codes.Unavailable:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_all_event_with_type").
				Msg("Service is down")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    "Not found",
				Data:       nil,
			}
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    "Forbidden resource",
				Data:       nil,
			}
		default:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_all_event_with_type").
				Msg("Unhandled error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return res, nil
}

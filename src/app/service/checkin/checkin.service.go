package checkin

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	client proto.CheckinServiceClient
}

func NewService(client proto.CheckinServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) CheckinVerify(userid string, eventType int) (*proto.CheckinVerifyResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.CheckinVerify(ctx, &proto.CheckinVerifyRequest{
		Id:        userid,
		EventType: int32(eventType),
	})

	if errRes != nil {
		st, ok := status.FromError(errRes)
		if !ok {
			log.Error().
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinverify").
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
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinverify").
				Msg("Error while connecting to service")

			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}
		case codes.Internal:
			log.Error().
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinverify").
				Msg("Internal Server Error")

			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinverify").
				Msg("Unhandled Error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return &proto.CheckinVerifyResponse{
		CheckinToken: res.CheckinToken,
		CheckinType:  res.CheckinType,
	}, nil
}

func (s *Service) CheckinConfirm(token string) (*proto.CheckinConfirmResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.CheckinConfirm(ctx, &proto.CheckinConfirmRequest{
		Token: token,
	})

	if errRes != nil {
		st, ok := status.FromError(errRes)
		if !ok {
			log.Error().
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinconfirm").
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
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinconfirm").
				Msg("Error while connecting to service")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}

		case codes.InvalidArgument:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid Checkin type",
				Data:       nil,
			}
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    "Invalid Token",
				Data:       nil,
			}
		default:
			log.Error().
				Err(errRes).
				Str("service", "checkin").
				Str("module", "checkinverify").
				Msg("Unhandled Error (Possibly invalid user id)")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return res, nil
}

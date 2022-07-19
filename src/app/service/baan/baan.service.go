package baan

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
	client proto.BaanServiceClient
}

func NewService(client proto.BaanServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) FindAll() ([]*proto.Baan, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Print("Find all baan")

	res, err := s.client.FindAllBaan(ctx, &proto.FindAllBaanRequest{})
	if err != nil {
		log.Error().Err(err).
			Str("service", "baan").
			Str("module", "find all").
			Msg("Error while find all baan")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
		}
	}

	return res.Baans, nil
}

func (s *Service) FindOne(id string) (*proto.Baan, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindOneBaan(ctx, &proto.FindOneBaanRequest{Id: id})
	if err != nil {

		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				log.Error().Err(err).
					Str("service", "baan").
					Str("module", "find one").
					Str("baan_id", id).
					Msg("Not found baan")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    "Baan not found",
				}
			default:
				log.Error().
					Err(err).
					Str("service", "user").
					Str("module", "findOne").
					Str("baan_id", id).
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
			Str("service", "user").
			Str("module", "findOne").
			Str("baan_id", id).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Baan, nil
}

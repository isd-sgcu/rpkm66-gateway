package file

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
	client proto.FileServiceClient
}

func NewService(client proto.FileServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) UploadImage(file *dto.DecomposedFile) (string, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.UploadImage(ctx, &proto.UploadImageRequest{
		Filename: file.Filename,
		Data:     file.Data,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unavailable:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload image").
					Msg("Cannot connect to google cloud storage")
				return "", &dto.ResponseErr{
					StatusCode: http.StatusGatewayTimeout,
					Message:    "Connection timeout",
					Data:       nil,
				}

			default:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload image").
					Msg("Error while connecting to service")
				return "", &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "upload image").
			Msg("Error while connecting to service")

		return "", &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}

	}

	return res.Filename, nil
}

package file

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/constant/file"
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

func (s *Service) Upload(file *dto.DecomposedFile, userId string, tag file.Tag, fileType file.Type) (string, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Upload(ctx, &proto.UploadRequest{
		Filename: file.Filename,
		Data:     file.Data,
		UserId:   userId,
		Tag:      int32(tag),
		Type:     int32(fileType),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unavailable:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload").
					Msg("Something wrong")
				return "", &dto.ResponseErr{
					StatusCode: http.StatusGatewayTimeout,
					Message:    "Connection timeout",
					Data:       nil,
				}

			default:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload").
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
			Str("module", "upload").
			Msg("Error while connecting to service")

		return "", &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}

	}

	return res.Url, nil
}
